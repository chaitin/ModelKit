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
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
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
		inputMsg[0].MultiContent = []schema.ChatMessagePart{
			{
				Type: schema.ChatMessagePartTypeText,
				Text: "What's in the picture? Only answer me a word.",
			},
			{
				Type: schema.ChatMessagePartTypeImageURL,
				ImageURL: &schema.ChatMessageImageURL{
					URL:    imageURL,
					Detail: schema.ImageURLDetailAuto,
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
	u.Path = path.Join(u.Path, "/models")

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

func filterModelsByType(models []domain.ModelListItem, req *domain.ModelListReq) []domain.ModelListItem {
	raw := strings.ToLower(req.Type)
	p := strings.ToLower(req.Provider)
	switch raw {
	// 分析模型 排除 嵌入模型和重排模型
	case "analysis":
		filtered := make([]domain.ModelListItem, 0, len(models))
		for _, it := range models {
			if !isEmbeddingModel(it.Model, p) && !isRerankModel(it.Model) {
				filtered = append(filtered, it)
			}
		}
		return filtered
	// 分析模型-视觉模型 仅包含 视觉模型
	case "analysis-vl":
		filtered := make([]domain.ModelListItem, 0, len(models))
		for _, it := range models {
			if isVisionModel(it.Model, p) {
				filtered = append(filtered, it)
			}
		}
		return filtered
	// 聊天模型 排除 嵌入模型和重排模型
	case "chat", "llm":
		filtered := make([]domain.ModelListItem, 0, len(models))
		for _, it := range models {
			if !isEmbeddingModel(it.Model, p) && !isRerankModel(it.Model) {
				filtered = append(filtered, it)
			}
		}
		return filtered
	// 嵌入模型 仅包含 嵌入模型
	case "embedding":
		filtered := make([]domain.ModelListItem, 0, len(models))
		for _, it := range models {
			if isEmbeddingModel(it.Model, p) {
				filtered = append(filtered, it)
			}
		}
		return filtered
	// 重排模型 仅包含 重排模型
	case "reranker", "rerank":
		filtered := make([]domain.ModelListItem, 0, len(models))
		for _, it := range models {
			if isRerankModel(it.Model) {
				filtered = append(filtered, it)
			}
		}
		return filtered
	case "code", "coder":
		filtered := make([]domain.ModelListItem, 0, len(models))
		for _, it := range models {
			if isCodeModel(it.Model, p) {
				filtered = append(filtered, it)
			}
		}
		return filtered
	default:
		return models
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
