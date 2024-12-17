package models

import "time"

type ChatMessage struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Doc       string    `json:"doc"`
	Diff      string    `json:"diff"`
	CreatedAt time.Time `json:"created_at"`
}
