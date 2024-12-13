package routes

import (
	"net/http"

	"composer/internal/db"
	"composer/internal/models"

	"github.com/labstack/echo/v4"
)

func RegisterChatSessionRoutes(e *echo.Echo, database *db.Db) {
	e.POST("/chat-sessions", createChatSession(database))
	e.GET("/chat-sessions", listChatSessions(database))
	e.GET("/chat-sessions/:id", getChatSession(database))
	e.PUT("/chat-sessions/:id", updateChatSession(database))
	e.DELETE("/chat-sessions/:id", deleteChatSession(database))
}

func createChatSession(database *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := database.InsertChatSession(&models.ChatSession{}); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, nil)
	}
}

func listChatSessions(database *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		chatSessions, err := database.ListChatSessions()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, chatSessions)
	}
}

func getChatSession(database *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		chatSession, err := database.GetChatSession(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}

		return c.JSON(http.StatusOK, chatSession)
	}
}

func updateChatSession(database *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		var chatSession models.ChatSession
		if err := c.Bind(&chatSession); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		chatSession.ID = id
		if err := database.UpdateChatSession(&chatSession); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, chatSession)
	}
}

func deleteChatSession(database *db.Db) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		if err := database.DeleteChatSession(id); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
