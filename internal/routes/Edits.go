package routes

import (
	"composer/internal/db"
	"composer/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/llms"
)

func RegisterEditRoutes(e *echo.Echo) {
	e.POST("/api/chat-sessions/:id/edits", handleEdits)
}

func handleEdits(c echo.Context) error {
	sessionID := c.Param("id")
	database := c.Get("db").(*db.Db)
	aiModel := c.Get("llm").(llms.Model)

	previousArtifact, err := getPreviousArtifactVersion(database, sessionID, "")
	if err != nil {
		return err
	}

	editsSystemPrompt := fmt.Sprintf(`
	You are a helpful assistant. For the request made to you, please provide your
	response in the following format, where you are providing the contents of the artifact
	you are generating and your thought process in distinctly demarcated tags. Please ensure that
	you do not provide any content outside of these tags.

	You will be given an existing artifact, in addition, you will also be given the changes the user has requested. 
	The requested edits will be given to you in the form of text snippets from the artifact that the user has placed
	their comments against to explain the change they'd like you to make. You will have to rewrite / update / edits parts
	of the artifact based on this input.

	<artifact>%s</artifact>

	Please think though the various sections of the document and put in as much detail as possible. Be thorough and detailed.
	Please respond in the following format including edits in between the edit tags

	<edit>
		<textToReplace>The exact original text to replace. This may or may not be same as the user selected text</textToReplace>
		<replacementText> Text to replace the original text snippet with</replacementText>
	</edit>
	`, previousArtifact)

	req := []requestedEdits{}
	err = c.Bind(&req)
	if err != nil {
		return err
	}

	var content string

	for _, re := range req {
		content += fmt.Sprintf("\nSelected text snippet from artifact:%s\nUser Comment:%s\n-----", re.Text, re.Comment)
	}
	msg := models.ChatMessage{
		SessionID: sessionID,
		Role:      "human",
		Content:   content,
		CreatedAt: time.Now(),
	}

	err = database.InsertChatMessage(&msg)
	if err != nil {
		return err
	}

	messageToModel := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageType("system"), editsSystemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, content),
	}

	result, err := aiModel.GenerateContent(c.Request().Context(), messageToModel, llms.WithMaxTokens(8192), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {

	}))

	log.Printf("Edits response %s", result.Choices[0].Content)
}

type requestedEdits struct {
	Text    string `json:"text"`
	Comment string `json:"comment"`
}
