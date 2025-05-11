package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"resumoturbo-api/internal/openai"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	openaiAPI "github.com/sashabaranov/go-openai"
)

type ResumoRequest struct {
	Texto string `json:"texto"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️  Não foi possível carregar .env:", err)
	}

	// Rodar sem logs em produção
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ou "http://localhost:8081"
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rota de resumo
	r.POST("/resumir", func(c *gin.Context) {
		var req ResumoRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Texto == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Texto inválido"})
			return
		}

		resumo, topicos, perguntas, err := openai.ProcessarTexto(req.Texto)
		if err != nil {
			fmt.Println("❌ Erro ao processar IA:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar resumo com IA"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"resumo":    resumo,
			"topicos":   topicos,
			"perguntas": perguntas,
		})
	})

	// Rota de status com teste real da OpenAI
	r.GET("/status", func(c *gin.Context) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ok":    false,
				"error": "Chave da OpenAI não configurada",
			})
			return
		}

		client := openaiAPI.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openaiAPI.ChatCompletionRequest{
				Model: openaiAPI.GPT3Dot5Turbo,
				Messages: []openaiAPI.ChatCompletionMessage{
					{
						Role:    openaiAPI.ChatMessageRoleUser,
						Content: "Diga apenas: API online",
					},
				},
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ok":      false,
				"error":   "Falha ao se comunicar com a OpenAI",
				"detalhe": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":        true,
			"status":    "API online",
			"openai_ok": resp.Choices[0].Message.Content,
			"version":   "1.0.0",
		})
	})

	// Porta padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
