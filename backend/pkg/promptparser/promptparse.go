package promptparser

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/chaitin/ModelKit/backend/pkg/cvt"
)

type PromptParser interface {
	Parse(prompt string) (*Info, error)
}

type Info struct {
	Prompt          string
	ProgramLanguage string
	WorkMode        string
}

type Kind int

const (
	KindNormal Kind = iota + 1
	KindTask
)

func New(kind Kind) PromptParser {
	switch kind {
	case KindNormal:
		return &NormalParser{}
	case KindTask:
		return &TaskParse{}
	default:
		return nil
	}
}

type NormalParser struct {
}

// Parse implements PromptParser.
func (n *NormalParser) Parse(prompt string) (*Info, error) {
	// re := regexp.MustCompile(`You are a code completion assistant for (\w+). Complete the code, don't explain it. Only return code that should come next.(.*)`)
	re := regexp.MustCompile(`You are a code completion assistant for (\w+). Complete the code naturally, providing complete and meaningful code blocks when appropriate. Return only the code that should come next, without explanations.(.*)`)
	match := re.FindStringSubmatch(prompt)
	if len(match) < 3 {
		return nil, fmt.Errorf("invalid prompt")
	}
	return &Info{
		ProgramLanguage: match[1],
		Prompt:          match[2],
	}, nil
}

type TaskParse struct {
	Task        string `xml:"task"`
	Feedback    string `xml:"feedback"`
	UserMessage string `xml:"user_message"`
}

func (m *TaskParse) Parse(prompt string) (*Info, error) {
	var tp TaskParse
	prompt = "<root>" + prompt + "</root>"
	if err := xml.Unmarshal([]byte(prompt), &tp); err != nil {
		return nil, err
	}

	return &Info{
		Prompt: cvt.CanditionVar(func() (string, bool) {
			return tp.Task, tp.Task != ""
		}, func() (string, bool) {
			return tp.Feedback, tp.Feedback != ""
		}, func() (string, bool) {
			return tp.UserMessage, tp.UserMessage != ""
		}),
	}, nil
}
