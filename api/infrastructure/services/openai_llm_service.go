package infrastructure_services

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAILlmService struct {
	llm *openai.LLM
}

func NewOpenAILlmService() *OpenAILlmService {
	llm, err := openai.New(openai.WithModel("gpt-4o-mini"))
	if err != nil {
		log.Fatal(err.Error())
	}

	return &OpenAILlmService{
		llm: llm,
	}
}

func (l *OpenAILlmService) SinglePrompt(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, l.llm, prompt)
}
