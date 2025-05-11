# 🧠 ResumoTurbo API (Golang + OpenAI)

Backend leve e eficiente feito em **Go** que utiliza a **API da OpenAI** para gerar:

- Resumo automático de texto
- Tópicos principais
- Perguntas para estudo

---

## 🚀 Rotas disponíveis

### `POST /resumir`

Envia um texto e retorna resumo, tópicos e perguntas.

#### Requisição:

```json
{
  "texto": "Texto completo que será resumido."
}
