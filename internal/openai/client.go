package openai

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func ProcessarTexto(texto string) (resumo string, topicos []string, perguntas []string, err error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ OPENAI_API_KEY não está definida no .env")
		return "", nil, nil, fmt.Errorf("chave da OpenAI ausente")
	}

	client := openai.NewClient(apiKey)

	// Prompt para gerar resumo, tópicos e perguntas
	prompt := fmt.Sprintf(`Resuma o seguinte texto em até 5 parágrafos. Depois, liste de 3 a 5 tópicos principais e 3 perguntas para estudo.

Texto:
%s`, texto)

	// Envia a requisição à OpenAI
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Você é um assistente que cria resumos claros e didáticos."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
	if err != nil {
		fmt.Println("❌ Erro na chamada à OpenAI:", err)
		return "", nil, nil, err
	}

	content := resp.Choices[0].Message.Content

	// Separar blocos de texto por dois \n\n (parágrafos)
	parts := strings.Split(content, "\n\n")
	if len(parts) == 0 {
		return "", nil, nil, fmt.Errorf("resposta inesperada da OpenAI")
	}

	// Primeiro parágrafo = resumo
	resumo = parts[0]

	// Os demais contêm tópicos e perguntas
	for _, p := range parts[1:] {
		if strings.Contains(strings.ToLower(p), "tópicos") {
			for _, l := range strings.Split(p, "\n") {
				if strings.HasPrefix(l, "-") || strings.HasPrefix(l, "•") {
					topicos = append(topicos, strings.TrimSpace(l[1:]))
				}
			}
		}

		if strings.Contains(strings.ToLower(p), "perguntas") {
			for _, l := range strings.Split(p, "\n") {
				if strings.HasPrefix(strings.TrimSpace(l), "1.") || strings.HasPrefix(strings.TrimSpace(l), "2.") || strings.HasPrefix(strings.TrimSpace(l), "3.") {
					perguntas = append(perguntas, strings.TrimSpace(l))
				}
			}
		}
	}

	return resumo, topicos, perguntas, nil
}
