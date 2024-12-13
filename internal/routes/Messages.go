package routes

import (
	"composer/internal/db"
	"composer/internal/models"
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/llms"
)

func RegisterMessageRoutes(e *echo.Echo) {
	e.POST("/chat-sessions/:id/messages", createMessage)
}

func createMessage(c echo.Context) error {
	sessionID := c.Param("id")
	database := c.Get("db").(*db.Db)
	aiModel := c.Get("llm").(llms.Model)

	rb := requestBody{}
	err := c.Bind(&rb)
	if err != nil {
		return err
	}

	msg := models.ChatMessage{
		SessionID: sessionID,
		Role:      "human",
		Content:   rb.Content,
		CreatedAt: time.Now(),
	}

	err = database.InsertChatMessage(&msg)
	if err != nil {
		return err
	}

	history, err := database.ListChatMessages(sessionID)
	if err != nil {
		return err
	}

	messageToModel := []llms.MessageContent{}

	for _, m := range history {
		messageToModel = append(messageToModel, llms.TextParts(llms.ChatMessageType(m.Role), m.Content))
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	aiResponse, err := aiModel.GenerateContent(
		c.Request().Context(),
		messageToModel,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			w.Write(chunk)
			w.Flush()
			return nil
		}),
	)
	if err != nil {
		return err
	}

	err = database.InsertChatMessage(&models.ChatMessage{
		SessionID: sessionID,
		Role:      "ai",
		Content:   aiResponse.Choices[0].Content,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

type requestBody struct {
	Content string `json:"content"`
}
