package db

import "composer/internal/models"

func (d *Db) InsertChatMessage(msg *models.ChatMessage) error {
	query := `
	INSERT INTO chat_messages (session_id, role, content, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	err := d.conn.QueryRow(query, msg.SessionID, msg.Role, msg.Content, msg.CreatedAt).Scan(&msg.ID)
	return err
}

func (d *Db) UpdateChatMessage(msg *models.ChatMessage) error {
	query := `
	UPDATE chat_messages 
	SET session_id = $1, role = $2, content = $3 
	WHERE id = $4`

	_, err := d.conn.Exec(query, msg.SessionID, msg.Role, msg.Content, msg.ID)
	return err
}

func (d *Db) DeleteChatMessage(id string) error {
	query := `DELETE FROM chat_messages WHERE id = $1`
	_, err := d.conn.Exec(query, id)
	return err
}

func (d *Db) GetChatMessage(id string) (*models.ChatMessage, error) {
	query := `
	SELECT id, session_id, role, content, created_at 
	FROM chat_messages 
	WHERE id = $1`

	msg := &models.ChatMessage{}
	err := d.conn.QueryRow(query, id).Scan(
		&msg.ID,
		&msg.SessionID,
		&msg.Role,
		&msg.Content,
		&msg.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (d *Db) ListChatMessages(sessionID string) ([]*models.ChatMessage, error) {
	query := `
	SELECT id, session_id, role, content, created_at 
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
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
