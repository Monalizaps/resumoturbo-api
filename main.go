package main

import (
	"net/http"
	"os"
	"resumoturbo-api/internal/openai"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type ResumoRequest struct {
	Texto string `json:"texto"`
}

func main() {
	godotenv.Load()

	r := gin.Default()
	r.POST("/resumir", func(c *gin.Context) {
		var req ResumoRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Texto == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Texto inválido"})
			return
		}

		resumo, topicos, perguntas, err := openai.ProcessarTexto(req.Texto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar IA"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"resumo":    resumo,
			"topicos":   topicos,
			"perguntas": perguntas,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
