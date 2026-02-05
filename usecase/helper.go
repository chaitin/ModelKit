package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	generativeGenai "github.com/google/generative-ai-go/genai"
	"github.com/ollama/ollama/api"
	"google.golang.org/api/option"
	"google.golang.org/genai"

	"github.com/chaitin/ModelKit/v2/consts"
	"github.com/chaitin/ModelKit/v2/domain"
	"github.com/chaitin/ModelKit/v2/pkg/request"
	"github.com/chaitin/ModelKit/v2/utils"
)

// 以下是辅助函数，用于处理模型列表和检查相关的功能
func ollamaListModel(baseURL string, httpClient *http.Client, apiHeader string) (*domain.ModelListResp, error) {
	// get from ollama http://10.10.16.24:11434/api/tags
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/api/tags"
	client := request.NewClient(u.Scheme, u.Host, httpClient.Timeout, request.WithClient(httpClient))

	h := request.Header{}
	if apiHeader != "" {
		headers := request.GetHeaderMap(apiHeader)
		maps.Copy(h, headers)
	}
	return request.Get[domain.ModelListResp](client, u.Path, request.WithHeader(h))
}

func (m *ModelKit) getChatModelGenerateChat(ctx context.Context, provider consts.ModelProvider, modelType consts.ModelType, baseURL string, req *domain.CheckModelReq) (string, error) {
	chatModel, err := m.GetChatModel(ctx, &domain.ModelMetadata{
		Provider:   provider,
		ModelName:  req.Model,
		APIKey:     req.APIKey,
		APIHeader:  req.APIHeader,
		BaseURL:    baseURL,
		APIVersion: req.APIVersion,
		ModelType:  modelType,
	})
	if err != nil {
		return "", err
	}
	checkImage := req.Param != nil && req.Param.SupportImages

	// gemini 多模态检测
	if provider == consts.ModelProviderGemini && checkImage {
		resp, err := m.geminiImageCheck(ctx, req)
		if err != nil {
			return "", err
		}
		if !strings.Contains(strings.ToLower(resp), "dog") {
			return "", fmt.Errorf("this model not support image input")
		}
		return resp, nil
	}

	genResp, err := chatModel.Generate(ctx, getInputMsg(req))

	// 非流式生成失败，尝试流式生成
	if err != nil || genResp.Content == "" {
		if m.logger != nil {
			m.logger.Info("Generate chat failed", slog.Any("error", err))
		} else {
			log.Printf("Generate chat failed, err: %v", err)
		}

		streamRes, streamErr := streamCheck(ctx, &chatModel, req)
		if streamErr != nil {
			if m.logger != nil {
				m.logger.Info("Stream chat failed", slog.Any("error", streamErr))
			} else {
				log.Printf("Stream chat failed, err: %v", streamErr)
			}
			return "", err
		}
		return streamRes, nil
	}

	if checkImage && !strings.Contains(strings.ToLower(genResp.Content), "dog") {
		return "", fmt.Errorf("this model not support image input")
	}
	return genResp.Content, nil
}

func getInputMsg(req *domain.CheckModelReq) []*schema.Message {
	var inputMsg []*schema.Message
	checkImage := req.Param != nil && req.Param.SupportImages
	if checkImage {
		imageURL := consts.ImageBase64
		// ollama 非openai兼容
		if req.Provider == string(consts.ModelProviderOllama) && !strings.HasSuffix(req.BaseURL, "/v1") {
			_, currentFile, _, _ := runtime.Caller(0)
			currentDir := filepath.Dir(currentFile)
			imagePath := filepath.Join(currentDir, "assets", "image.png")
			image, _ := os.ReadFile(imagePath)
			imageURL = string(image)
		}
		inputMsg = []*schema.Message{
			schema.UserMessage(""),
		}
		inputMsg[0].UserInputMultiContent = []schema.MessageInputPart{
			{
				Type: schema.ChatMessagePartTypeText,
				Text: "What's in the picture? Only answer me a word.",
			},
			{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{URL: &imageURL},
					Detail:            schema.ImageURLDetailAuto,
				},
			},
		}
	} else {
		inputMsg = []*schema.Message{
			schema.SystemMessage("You are a helpful assistant."),
			schema.UserMessage("hi"),
		}
	}
	return inputMsg
}

func streamCheck(ctx context.Context, chatModel *model.BaseChatModel, req *domain.CheckModelReq) (string, error) {
	var res string

	streamResult, err := (*chatModel).Stream(ctx, getInputMsg(req))
	if err != nil {
		return "", err
	}

	for {
		chunk, err := streamResult.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		// 响应片段处理
		res += chunk.Content
	}
	return res, nil
}

// baseURL的host换成host.docker.internal
func baseURLReplaceHost(inputURL string) (string, error) {
	rawURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	var hostAddress string
	if runtime.GOOS == "linux" {
		hostAddress = consts.LinuxHost
	} else {
		hostAddress = consts.MacWinHost
	}

	if rawURL.Hostname() != hostAddress {
		if rawURL.Port() != "" {
			rawURL.Host = hostAddress + ":" + rawURL.Port()
		} else {
			rawURL.Host = hostAddress
		}
	}
	return rawURL.String(), nil
}

func baseURLReplaceSlash(inputURL string) (string, error) {
	rawURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	// 去掉末尾的/
	rawURL.Path = strings.TrimSuffix(rawURL.Path, "/")
	return rawURL.String(), nil
}

// reqModelListApi 获取OpenAI兼容API的模型列表
// 使用泛型和接口抽象来支持不同供应商的响应格式
func reqModelListApi[T domain.ModelResponseParser](req *domain.ModelListReq, httpClient *http.Client, responseType T) ([]domain.ModelListItem, error) {
	u, err := url.Parse(req.BaseURL)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(req.BaseURL, "#") {
		u.Fragment = ""
		// 使用原始路径（去掉结尾#），不再追加 /models
	} else {
		u.Path = path.Join(u.Path, "/models")
	}

	client := request.NewClient(u.Scheme, u.Host, httpClient.Timeout, request.WithClient(httpClient))
	query, err := utils.GetQuery(req)
	if err != nil {
		return nil, err
	}
	resp, err := request.Get[T](
		client, u.Path,
		request.WithHeader(
			request.Header{
				"Authorization": fmt.Sprintf("Bearer %s", req.APIKey),
			},
		),
		request.WithQuery(query),
	)
	if err != nil {
		return nil, err
	}

	return (*resp).ParseModels(), nil
}

func generateBaseURLFixSuggestion(errContent string, baseURL string, provider consts.ModelProvider) string {
	var is404, isLocal, hasPath, isOther, isEndWithSlash bool
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	if strings.HasSuffix(baseURL, "/") {
		isEndWithSlash = true
	}
	if strings.Contains(errContent, "404") || strings.Contains(errContent, "connection refused") {
		is404 = true
	}
	if strings.Contains(parsedURL.Host, consts.LocalHost) || strings.Contains(parsedURL.Host, consts.LocalIP) {
		isLocal = true
	}
	if parsedURL.Path != "" {
		hasPath = true
	}
	isOther = provider == consts.ModelProviderOther

	var errType consts.AddModelBaseURLErrType
	if strings.Contains(baseURL, "chat/completions") {
		errType = consts.AddModelBaseURLErrTypeChatCompletions
	} else if isEndWithSlash {
		errType = consts.AddModelBaseURLErrTypeSlash
	} else if is404 && isLocal { // 404 且是本地地址，建议使用宿主机主机名
		errType = consts.AddModelBaseURLErrTypeHost
	} else if !isLocal && !hasPath && isOther {
		// 不是本地地址，且没有path，建议在API地址末尾添加/v1
		errType = consts.AddModelBaseURLErrTypeV1Path
	} else {
		return ""
	}

	switch errType {
	case consts.AddModelBaseURLErrTypeHost:
		fixedURL, err := baseURLReplaceHost(baseURL)
		if err != nil {
			return ""
		}
		return "建议在API地址使用宿主机主机名: " + fixedURL
	case consts.AddModelBaseURLErrTypeSlash:
		fixedURL, err := baseURLReplaceSlash(baseURL)
		if err != nil {
			return ""
		}
		return "请去掉API地址末尾的/: " + fixedURL
	case consts.AddModelBaseURLErrTypeChatCompletions:
		return "请去掉/chat/completions路径"
	default:
		return ""
	}
}

func (m *ModelKit) geminiImageCheck(ctx context.Context, req *domain.CheckModelReq) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: req.APIKey,
	})
	if err != nil {
		return "", err
	}

	// 获取当前文件所在目录
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	imagePath := filepath.Join(currentDir, "assets", "image.png")
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/png"),
		genai.NewPartFromText("What's in the picture? Only answer me a word.If you don't support image input, reply no."),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		req.Model,
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

func FilterModelsByType(models []domain.ModelListItem, req *domain.ModelListReq) []domain.ModelListItem {
	p := strings.ToLower(req.Provider)
	t := normalizeType(strings.ToLower(req.Type))
	pred := modelPredicate(t, p)
	if pred == nil {
		return models
	}
	filtered := make([]domain.ModelListItem, 0, len(models))
	for _, it := range models {
		if pred(it.Model) {
			filtered = append(filtered, it)
		}
	}
	return filtered
}

func normalizeType(t string) string {
	switch t {
	case "chat", "llm":
		return "chat"
	case "reranker", "rerank":
		return "rerank"
	case "code", "coder":
		return "code"
	default:
		return t
	}
}

func modelPredicate(t, provider string) func(string) bool {
	switch t {
	case "analysis":
		return func(m string) bool { return !isEmbeddingModel(m, provider) && !isRerankModel(m) }
	case "analysis-vl":
		return func(m string) bool { return isVisionModel(m, provider) }
	case "chat":
		return func(m string) bool { return !isEmbeddingModel(m, provider) && !isRerankModel(m) }
	case "embedding":
		return func(m string) bool { return isEmbeddingModel(m, provider) }
	case "rerank":
		return func(m string) bool { return isRerankModel(m) }
	case "code":
		return func(m string) bool { return isCodeModel(m, provider) }
	default:
		return nil
	}
}

func getLowerBaseModelName(id string) string {
	parts := strings.Split(id, "/")
	return strings.ToLower(parts[len(parts)-1])
}

func isRerankModel(modelID string) bool {
	if modelID == "" {
		return false
	}
	mid := getLowerBaseModelName(modelID)
	re := regexp.MustCompile(`(?i)(?:rerank|re-rank|re-ranker|re-ranking|retrieval|retriever)`)
	return re.MatchString(mid)
}

func isEmbeddingModel(modelID, provider string) bool {
	if modelID == "" {
		return false
	}
	if isRerankModel(modelID) {
		return false
	}
	mid := getLowerBaseModelName(modelID)
	if provider == "anthropic" {
		return false
	}
	if provider == "doubao" || strings.Contains(mid, "doubao") {
		re := regexp.MustCompile(`(?i)(?:^text-|embed|bge-|e5-|LLM2Vec|retrieval|uae-|gte-|jina-clip|jina-embeddings|voyage-)`)
		return re.MatchString(mid)
	}
	re := regexp.MustCompile(`(?i)(?:^text-|embed|bge-|e5-|LLM2Vec|retrieval|uae-|gte-|jina-clip|jina-embeddings|voyage-)`)
	return re.MatchString(mid)
}

func isCodeModel(modelID, provider string) bool {
	if modelID == "" {
		return false
	}
	if isEmbeddingModel(modelID, provider) || isRerankModel(modelID) {
		return false
	}
	mid := getLowerBaseModelName(modelID)
	re := regexp.MustCompile(`(?i)(?:^o3$|.*(code|claude\s+sonnet|claude\s+opus|gpt-4\.1|gpt-4o|gpt-5|gemini[\s-]+2\.5|o4-mini|kimi-k2).*)`)
	return re.MatchString(mid)
}

var visionModels = []string{
	`chatgpt-4o(?:-[\w-]+)?`,
	`claude-3`,
	`claude-opus-4`,
	`claude-sonnet-4`,
	`deepseek-vl(?:[\w-]+)?`,
	`doubao-seed-1[.-]6(?:-[\w-]+)?`,
	`gemini-1\.5`,
	`gemini-2\.0`,
	`gemini-2\.5`,
	`gemini-exp`,
	`gemma-3(?:-[\w-]+)`,
	`gemma3(?:[-:\w]+)?`,
	`glm-4(?:\.\d+)?v(?:-[\w-]+)?`,
	`gpt-4(?:-[\w-]+)`,
	`gpt-4.1(?:-[\w-]+)?`,
	`gpt-4.5(?:-[\w-]+)`,
	`gpt-4o(?:-[\w-]+)?`,
	`gpt-5(?:-[\w-]+)?`,
	`grok-4(?:-[\w-]+)?`,
	`grok-vision-beta`,
	`internvl2`,
	`kimi-latest`,
	`kimi-thinking-preview`,
	`kimi-vl-a3b-thinking(?:-[\w-]+)?`,
	`llama-4(?:-[\w-]+)?`,
	`llama-guard-4(?:-[\w-]+)?`,
	`llava`,
	`minicpm`,
	`moondream`,
	`o1(?:-[\w-]+)?`,
	`o3(?:-[\w-]+)?`,
	`o4(?:-[\w-]+)?`,
	`pixtral`,
	`qvq`,
	`qwen-vl`,
	`qwen2-vl`,
	`qwen2.5-omni`,
	`qwen2.5-vl`,
	`step-1o(?:.*vision)?`,
	`step-1v(?:-[\w-]+)?`,
	`vision`,
}

var notVisionModels = []string{
	`AIDC-AI/Marco-o1`,
	`gpt-4-\d+-preview`,
	`gpt-4-turbo-preview`,
	`gpt-4-32k`,
	`gpt-4-\d+`,
	`o1-mini`,
	`o3-mini`,
	`o1-preview`,
}

func isVisionModel(modelID, provider string) bool {
	if modelID == "" {
		return false
	}
	if isEmbeddingModel(modelID, provider) || isRerankModel(modelID) {
		return false
	}
	mid := getLowerBaseModelName(modelID)
	not := strings.Join(notVisionModels, "|")
	yes := strings.Join(visionModels, "|")
	yesRe := regexp.MustCompile(`(?i)\b(?:` + yes + `)\b`)
	if !yesRe.MatchString(mid) {
		return false
	}
	notRe := regexp.MustCompile(`(?i)\b(?:` + not + `)\b`)
	return !notRe.MatchString(mid)
}

func (m *ModelKit) checkChatModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	checkResp := &domain.CheckModelResp{}
	req.BaseURL = strings.TrimSuffix(req.BaseURL, "#")
	provider := consts.ParseModelProvider(req.Provider)
	modelType := consts.ParseModelType(req.Type)

	resp, err := m.getChatModelGenerateChat(ctx, provider, modelType, req.BaseURL, req)
	if err != nil && (provider == consts.ModelProviderOther || provider == consts.ModelProviderOllama || provider == consts.ModelProviderAzureOpenAI) {
		msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
		if msg == "" {
			checkResp.Error = err.Error()
		} else {
			checkResp.Error = msg
		}
		return checkResp, nil
	}
	if err != nil && provider != consts.ModelProviderOther {
		errorMsg := strings.ToLower(err.Error())
		for _, keyword := range consts.ApiKeyBalanceKeyWords {
			if strings.Contains(errorMsg, keyword) {
				checkResp.Error = "API Key余额不足"
				return checkResp, nil
			}
		}
		checkResp.Error = err.Error()
		return checkResp, nil
	}
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}
	if resp == "" {
		checkResp.Error = "生成内容失败"
		return checkResp, nil
	}
	checkResp.Content = resp
	return checkResp, nil
}

func (m *ModelKit) checkEmbeddingModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	checkResp := &domain.CheckModelResp{}
	provider := consts.ParseModelProvider(req.Provider)
	encodingFormat := "float"

	embedder, err := m.GetEmbedder(ctx, &domain.ModelMetadata{
		Provider:  provider,
		ModelName: req.Model,
		BaseURL:   req.BaseURL,
		APIKey:    req.APIKey,
		EmbedderParam: domain.EmbedderParam{
			EncodingFormat: &encodingFormat,
		},
	})
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}

	embResp, err := m.UseEmbedder(ctx, embedder, []string{"ModelKit 一个轻量级工具库，提供 AI 模型发现与 API 密钥验证功能，助你快速集成各大模型供应商能力。"})
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}
	if embResp == nil || len(embResp.Embeddings) == 0 {
		checkResp.Error = "empty embeddings"
		return checkResp, nil
	}
	dim := len(embResp.Embeddings[0].Embedding)
	checkResp.Content = fmt.Sprintf("dim is : %d", dim)
	return checkResp, nil
}

func (m *ModelKit) checkRerankModel(ctx context.Context, req *domain.CheckModelReq) (*domain.CheckModelResp, error) {
	checkResp := &domain.CheckModelResp{}
	provider := consts.ParseModelProvider(req.Provider)

	reranker, err := m.GetReranker(ctx, &domain.ModelMetadata{
		Provider:  provider,
		ModelName: req.Model,
		BaseURL:   req.BaseURL,
		APIKey:    req.APIKey,
	})
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}

	res, err := reranker.Rerank(ctx, domain.RerankRequest{
		Query:           "动物",
		Documents:       []string{"鸡", "火", "面", "火鸡面", "火鸡"},
		ReturnDocuments: true,
	})
	if err != nil {
		checkResp.Error = err.Error()
		return checkResp, nil
	}

	docs := make([]string, 0, len(res.Results))
	for _, r := range res.Results {
		if r.Document != "" {
			docs = append(docs, r.Document)
		}
	}
	checkResp.Content = strings.Join(docs, "\n")
	return checkResp, nil
}

func (m *ModelKit) listStaticProvider(req *domain.ModelListReq, provider consts.ModelProvider) (*domain.ModelListResp, error) {
	models := domain.From(domain.ModelProviders[provider])
	filtered := FilterModelsByType(models, req)
	return &domain.ModelListResp{Models: filtered}, nil
}

func (m *ModelKit) listGemini(ctx context.Context, req *domain.ModelListReq) (*domain.ModelListResp, error) {
	client, err := generativeGenai.NewClient(ctx, option.WithAPIKey(req.APIKey))
	if err != nil {
		return &domain.ModelListResp{Error: err.Error()}, nil
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			if m.logger != nil {
				m.logger.Error("Failed to close gemini client: %v", slog.Any("error", closeErr))
			} else {
				log.Printf("Failed to close gemini client: %v", closeErr)
			}
		}
	}()

	modelsList := make([]domain.ModelListItem, 0)
	modelsIter := client.ListModels(ctx)
	for {
		model, err := modelsIter.Next()
		if err != nil {
			break
		}
		if !slices.Contains(model.SupportedGenerationMethods, "generateContent") {
			continue
		}
		if !strings.Contains(model.Name, "gemini") {
			continue
		}
		name, _ := strings.CutPrefix(model.Name, "models/")
		modelsList = append(modelsList, domain.ModelListItem{Model: name})
	}
	if len(modelsList) == 0 {
		return &domain.ModelListResp{Error: fmt.Errorf("获取Gemini模型列表失败: 未找到可用模型").Error()}, nil
	}
	filtered := FilterModelsByType(modelsList, req)
	return &domain.ModelListResp{Models: filtered}, nil
}

func (m *ModelKit) listGithub(req *domain.ModelListReq, httpClient *http.Client) (*domain.ModelListResp, error) {
	models, err := reqModelListApi(req, httpClient, &domain.GithubResp{})
	if err != nil {
		return &domain.ModelListResp{Error: err.Error()}, nil
	}
	filtered := FilterModelsByType(models, req)
	return &domain.ModelListResp{Models: filtered}, nil
}

func (m *ModelKit) listOllama(req *domain.ModelListReq, httpClient *http.Client) (*domain.ModelListResp, error) {
	var modelListResp domain.ModelListResp
	var err error
	if strings.HasSuffix(req.BaseURL, "/v1") {
		var models []domain.ModelListItem
		models, err = reqModelListApi(req, httpClient, &domain.OpenAIResp{})
		if err == nil {
			modelListResp.Models = FilterModelsByType(models, req)
		}
	} else {
		var resp *domain.ModelListResp
		resp, err = ollamaListModel(req.BaseURL, httpClient, req.APIHeader)
		if err == nil {
			modelListResp = *resp
			modelListResp.Models = FilterModelsByType(modelListResp.Models, req)
		}
	}
	if err != nil {
		provider := consts.ParseModelProvider(req.Provider)
		msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
		if msg == "" {
			return &domain.ModelListResp{Error: err.Error()}, nil
		}
		return &domain.ModelListResp{Error: msg}, nil
	}
	return &modelListResp, nil
}

func (m *ModelKit) listGPUStack(req *domain.ModelListReq, httpClient *http.Client) (*domain.ModelListResp, error) {
	provider := consts.ParseModelProvider(req.Provider)
	models, err := reqModelListApi(req, httpClient, &domain.GPUStackListModelResp{})
	if err != nil {
		if m.logger != nil {
			m.logger.Error("GPUStack list model failed", "error", err, "models: ", models)
		}
		msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
		if msg == "" {
			return &domain.ModelListResp{Error: err.Error()}, nil
		}
		return &domain.ModelListResp{Error: msg}, nil
	}
	filtered := FilterModelsByType(models, req)
	return &domain.ModelListResp{Models: filtered}, nil
}

func (m *ModelKit) listOpenAI(req *domain.ModelListReq, httpClient *http.Client, provider consts.ModelProvider) (*domain.ModelListResp, error) {
	models, err := reqModelListApi(req, httpClient, &domain.OpenAIResp{})
	if err != nil {
		if provider == consts.ModelProviderOllama {
			msg := generateBaseURLFixSuggestion(err.Error(), req.BaseURL, provider)
			if msg == "" {
				return &domain.ModelListResp{Error: err.Error()}, nil
			}
			return &domain.ModelListResp{Error: msg}, nil
		}
		return &domain.ModelListResp{Error: err.Error()}, nil
	}
	filtered := FilterModelsByType(models, req)
	return &domain.ModelListResp{Models: filtered}, nil
}

func buildOpenAIChatConfig(md *domain.ModelMetadata) *openai.ChatModelConfig {
	t := float32(0.0)
	if md.Temperature != nil {
		t = *md.Temperature
	}
	cfg := &openai.ChatModelConfig{
		APIKey:      md.APIKey,
		BaseURL:     md.BaseURL,
		Model:       string(md.ModelName),
		Temperature: &t,
	}
	if md.MaxTokens != nil {
		cfg.MaxTokens = md.MaxTokens
	}
	if md.TopP != nil {
		cfg.TopP = md.TopP
	}
	if len(md.Stop) > 0 {
		cfg.Stop = md.Stop
	}
	if md.PresencePenalty != nil {
		cfg.PresencePenalty = md.PresencePenalty
	}
	if md.FrequencyPenalty != nil {
		cfg.FrequencyPenalty = md.FrequencyPenalty
	}
	if md.ResponseFormat != nil {
		cfg.ResponseFormat = md.ResponseFormat
	}
	if md.Seed != nil {
		cfg.Seed = md.Seed
	}
	if md.LogitBias != nil {
		cfg.LogitBias = md.LogitBias
	}
	if md.Provider == consts.ModelProviderAzureOpenAI {
		cfg.ByAzure = true
		cfg.APIVersion = md.APIVersion
		if cfg.APIVersion == "" {
			cfg.APIVersion = "2024-10-21"
		}
		cfg.AzureModelMapperFunc = func(model string) string {
			return model
		}
	}
	if md.APIHeader != "" {
		hc := utils.GetHttpClientWithAPIHeaderMap(md.APIHeader)
		if hc != nil {
			cfg.HTTPClient = hc
		}
	}
	return cfg
}

func newDeepseekChatModel(ctx context.Context, md *domain.ModelMetadata) (model.BaseChatModel, error) {
	t := float32(0.0)
	if md.Temperature != nil {
		t = *md.Temperature
	}
	cfg := &deepseek.ChatModelConfig{
		BaseURL:     md.BaseURL,
		APIKey:      md.APIKey,
		Model:       md.ModelName,
		Temperature: t,
	}
	if md.MaxTokens != nil {
		cfg.MaxTokens = *md.MaxTokens
	}
	if md.TopP != nil {
		cfg.TopP = *md.TopP
	}
	if len(md.Stop) > 0 {
		cfg.Stop = md.Stop
	}
	if md.PresencePenalty != nil {
		cfg.PresencePenalty = *md.PresencePenalty
	}
	if md.FrequencyPenalty != nil {
		cfg.FrequencyPenalty = *md.FrequencyPenalty
	}
	return deepseek.NewChatModel(ctx, cfg)
}

func newGeminiChatModel(ctx context.Context, md *domain.ModelMetadata) (model.BaseChatModel, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: md.APIKey})
	if err != nil {
		return nil, err
	}
	cfg := &gemini.Config{
		Client: client,
		Model:  md.ModelName,
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
			ThinkingBudget:  nil,
		},
	}
	if md.MaxTokens != nil {
		cfg.MaxTokens = md.MaxTokens
	}
	if md.Temperature != nil {
		cfg.Temperature = md.Temperature
	}
	if md.TopP != nil {
		cfg.TopP = md.TopP
	}
	return gemini.NewChatModel(ctx, cfg)
}

func newOllamaChatModel(ctx context.Context, md *domain.ModelMetadata) (model.BaseChatModel, error) {
	if strings.HasSuffix(md.BaseURL, "/v1") {
		cfg := buildOpenAIChatConfig(md)
		return openai.NewChatModel(ctx, cfg)
	}
	t := float32(0.0)
	if md.Temperature != nil {
		t = *md.Temperature
	}
	baseUrl, err := utils.URLRemovePath(md.BaseURL)
	if err != nil {
		return nil, err
	}
	opts := &api.Options{Temperature: t}
	if md.TopP != nil {
		opts.TopP = *md.TopP
	}
	if len(md.Stop) > 0 {
		opts.Stop = md.Stop
	}
	if md.PresencePenalty != nil {
		opts.PresencePenalty = *md.PresencePenalty
	}
	if md.FrequencyPenalty != nil {
		opts.FrequencyPenalty = *md.FrequencyPenalty
	}
	if md.Seed != nil {
		opts.Seed = *md.Seed
	}
	return ollama.NewChatModel(ctx, &ollama.ChatModelConfig{BaseURL: baseUrl, Model: string(md.ModelName), Options: opts})
}
