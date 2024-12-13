package models

import "time"

type Document struct {
	ID             string    `json:"id"`
	Contents       string    `json:"contents"`
	LastModifiedBy string    `json:"last_modified_by"`
	ChatMessageID  string    `json:"chat_message_id"`
	CreatedAt      time.Time `json:"created_at"`
}
