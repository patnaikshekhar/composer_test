package db

import "composer/internal/models"

func (d *Db) InsertChatMessage(msg *models.ChatMessage) error {
	query := `
	INSERT INTO chat_messages (session_id, role, content, doc, diff, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	err := d.conn.QueryRow(query, msg.SessionID, msg.Role, msg.Content, msg.Doc, msg.Diff, msg.CreatedAt).Scan(&msg.ID)
	return err
}

func (d *Db) ListChatMessages(sessionID string) ([]*models.ChatMessage, error) {
	query := `
	SELECT id, session_id, role, content, doc, diff, created_at 
	FROM chat_messages 
	WHERE session_id = $1 
	ORDER BY created_at`

	rows, err := d.conn.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		msg := &models.ChatMessage{}
		err := rows.Scan(
			&msg.ID,
			&msg.SessionID,
			&msg.Role,
			&msg.Content,
			&msg.Doc,
			&msg.Diff,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
