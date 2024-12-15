package routes

import (
	"composer/internal/db"
	"composer/internal/models"
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/llms"
)

func RegisterMessageRoutes(e *echo.Echo) {
	e.POST("/api/chat-sessions/:id/messages", createMessage)
}

const systemPrompt = `
You are a helpful assistant. For the request made to you, please provide your
response in the following format, where you are providing the contents of the artifact
you are generating and your thought process in distinctly demarcated tags. Please ensure that
you do not provide any content outside of these tags.

If you are asked to generate a document, please think though the various sections of the document and put in as
much detail as possible. Be thorough and detailed.
<artifact>
# Markdown of the document

## Section 1
...
</artifact>
<explanation>
Users first need to install the pre-requisites because...
</explanation>
	`

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

	messageToModel := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageType("system"), systemPrompt),
	}

	for _, m := range history {
		messageToModel = append(messageToModel, llms.TextParts(llms.ChatMessageType(m.Role), m.Content))
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	streamMessage := &UserChatMessageResponse{
		Message:  "",
		Artifact: "",
	}

	inArtifact := false
	inExplanation := false
	collectedChunks := ""

	result, err := aiModel.GenerateContent(c.Request().Context(), messageToModel, llms.WithMaxTokens(8192), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		log.Printf("Chunk is %s", chunk)
		collectedChunks += string(chunk)

		if strings.Contains(collectedChunks, "<artifact>") {
			inArtifact = true
			collectedChunks = strings.Replace(collectedChunks, "<artifact>", "", -1)
			streamMessage.Artifact += collectedChunks
			return nil
		}

		if strings.Contains(collectedChunks, "</artifact>") {
			inArtifact = false
			frags := strings.Split(string(chunk), "</artifact>")
			streamMessage.Artifact += frags[0]
			log.Printf("Ending artifact with message %s", streamMessage.Artifact)
			err := json.NewEncoder(w).Encode(streamMessage)
			if err != nil {
				return err
			}
			w.Write([]byte("\r\n"))
			w.Flush()
			collectedChunks = ""
			if len(frags) > 1 {
				collectedChunks = frags[1]
			}
			return nil
		}

		if inArtifact {
			streamMessage.Artifact += string(chunk)
			streamMessage.Artifact = strings.Replace(streamMessage.Artifact, "</artifact", "", -1)
			err := json.NewEncoder(w).Encode(streamMessage)
			if err != nil {
				return err
			}
			w.Write([]byte("\r\n"))
			w.Flush()
		}

		if strings.Contains(collectedChunks, "<explanation>") {
			inExplanation = true
			collectedChunks = strings.Replace(collectedChunks, "<explanation>", "", -1)
			streamMessage.Message += collectedChunks
			return nil
		}

		if strings.Contains(collectedChunks, "</explanation>") {
			inExplanation = false
			frags := strings.Split(string(chunk), "</explanation>")
			streamMessage.Message += frags[0]
			err := json.NewEncoder(w).Encode(streamMessage)
			if err != nil {
				return err
			}
			w.Write([]byte("\r\n"))
			w.Flush()
			if len(frags) > 1 {
				collectedChunks = frags[1]
			}
			return nil
		}

		if inExplanation {
			streamMessage.Message += string(chunk)
			// streamMessage.Message = strings.Replace(streamMessage.Message, "</explanation", "", -1)
			err := json.NewEncoder(w).Encode(streamMessage)
			if err != nil {
				return err
			}
			w.Write([]byte("\r\n"))
			w.Flush()
		}

		return nil
	}))
	if err != nil {
		return err
	}

	err = database.InsertChatMessage(&models.ChatMessage{
		SessionID: sessionID,
		Role:      "ai",
		Content:   streamMessage.Message,
		Doc:       streamMessage.Artifact,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	log.Printf("response from LLM is %+v", result.Choices[0])
	return nil
}

type requestBody struct {
	Content string `json:"content"`
}

type UserChatMessageResponse struct {
	Message  string `json:"message"`
	Artifact string `json:"artifact,omitempty"`
}
