package db

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"fmt"
)
var DB *sql.DB
func Init() {
	var err error
	DB, err = sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal("Error while database openinig: ", err)
		panic(err)
	}
	defer DB.Close()
	migrate()
}


func migrate() error {
	if _, err := DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := []string{

		`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL,
			name TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			logged_in INTEGER NOT NULL DEFAULT 0,
			conn_type TEXT NOT NULL
		);
		`,

		// rooms
		`
		CREATE TABLE IF NOT EXISTS rooms (
			id INTEGER NOT NULL,
			owner_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// room users (many-to-many)
		`
		CREATE TABLE IF NOT EXISTS room_users (
			room_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (room_id, user_id),
			FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// room admins
		`
		CREATE TABLE IF NOT EXISTS room_admins (
			room_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (room_id, user_id),
			FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// chats
		`
		CREATE TABLE IF NOT EXISTS chats (
			id INTEGER NOT NULL,
			user1_id INTEGER NOT NULL,
			user2_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (user1_id, user2_id),
			FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// messages
		`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			room_id INTEGER,
			chat_id INTEGER,
			body TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
			FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
		);
		`,

		// connection requests
		`
		CREATE TABLE IF NOT EXISTS connection_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user_id INTEGER NOT NULL,
			to_user_id INTEGER NOT NULL,
			message TEXT,
			status INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,
	}

	for i, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}

	return tx.Commit()
}

