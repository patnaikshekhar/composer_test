package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"composer/internal/models"
)

func (d *Db) InsertChatSession(chatSession *models.ChatSession) error {
	query := `INSERT INTO chat_sessions (title) VALUES (?)`
	log.Printf("Running query %s", query)
	result, err := d.conn.Exec(query, chatSession.Title)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	chatSession.ID = fmt.Sprintf("%d", id)
	return nil
}

func (d *Db) UpdateChatSession(chatSession *models.ChatSession) error {
	query := `UPDATE chat_sessions SET title = ? WHERE id = ?`
	_, err := d.conn.Exec(query, chatSession.Title, chatSession.ID)
	return err
}

func (d *Db) DeleteChatSession(id string) error {
	query := `DELETE FROM chat_sessions WHERE id = ?`
	_, err := d.conn.Exec(query, id)
	return err
}

func (d *Db) GetChatSession(id string) (*models.ChatSession, error) {
	query := `SELECT id, title, created_at FROM chat_sessions WHERE id = ?`
	row := d.conn.QueryRow(query, id)

	var chatSession models.ChatSession
	err := row.Scan(&chatSession.ID, &chatSession.Title, &chatSession.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("chat session not found")
		}
		return nil, err
	}

	return &chatSession, nil
}

func (d *Db) ListChatSessions() ([]models.ChatSession, error) {
	query := `SELECT id, title, created_at FROM chat_sessions`
	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatSessions []models.ChatSession
	for rows.Next() {
		var chatSession models.ChatSession
		err := rows.Scan(&chatSession.ID, &chatSession.Title, &chatSession.CreatedAt)
		if err != nil {
			return nil, err
		}
		chatSessions = append(chatSessions, chatSession)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chatSessions, nil
}
