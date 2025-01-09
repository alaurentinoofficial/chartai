package infrastructure_services

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type ClaudeLlmService struct {
	llm *anthropic.LLM
}

func NewClaudeLlmService() *ClaudeLlmService {
	llm, err := anthropic.New(anthropic.WithModel("claude-3-haiku-20240307"))
	if err != nil {
		log.Fatal(err.Error())
	}

	return &ClaudeLlmService{
		llm: llm,
	}
}

func (l *ClaudeLlmService) SinglePrompt(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, l.llm, prompt)
}
