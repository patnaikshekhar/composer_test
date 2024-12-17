package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	conn *sql.DB
}

func New(relationDBToUse, connectionString string) (*Db, error) {

	createTables := []string{
		`
	CREATE TABLE IF NOT EXISTS chat_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`,

		`
	CREATE TABLE IF NOT EXISTS chat_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT,
		doc TEXT,
		diff TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`,

		`
	CREATE TABLE IF NOT EXISTS documents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		contents TEXT NOT NULL,
		last_modified_by TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`,
	}

	conn, err := sql.Open(relationDBToUse, connectionString)
	if err != nil {
		return nil, err
	}

	for _, query := range createTables {
		_, err := conn.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return &Db{
		conn: conn,
	}, nil
}
