package openai

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func ProcessarTexto(texto string) (resumo string, topicos []string, perguntas []string, err error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	prompt := fmt.Sprintf(`Resuma o seguinte texto em até 5 parágrafos. Depois, liste de 3 a 5 tópicos principais e 3 perguntas para estudo:

Texto:
%s`, texto)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Você é um assistente que cria resumos claros e didáticos."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
	if err != nil {
		return "", nil, nil, err
	}

	content := resp.Choices[0].Message.Content

	// Separar conteúdo com base em marcadores
	parts := strings.Split(content, "\n\n")
	resumo = parts[0]

	for _, p := range parts[1:] {
		if strings.Contains(p, "Tópicos") || strings.Contains(p, "Principais") {
			lines := strings.Split(p, "\n")
			for _, l := range lines {
				if strings.HasPrefix(l, "-") || strings.HasPrefix(l, "•") {
					topicos = append(topicos, strings.TrimSpace(l[1:]))
				}
			}
		}
		if strings.Contains(p, "Perguntas") || strings.Contains(p, "Questões") {
			lines := strings.Split(p, "\n")
			for _, l := range lines {
				if strings.HasPrefix(l, "1.") || strings.HasPrefix(l, "2.") || strings.HasPrefix(l, "3.") {
					perguntas = append(perguntas, strings.TrimSpace(l))
				}
			}
		}
	}

	return resumo, topicos, perguntas, nil
}
