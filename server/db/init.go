package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DB"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Cannot ping database: ", err)
	}

	err = migrate()
	if err != nil {
		log.Fatal("Error while migrating the database: ", err)
	}
}

func migrate() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := []string{
		// users
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			logged_in BOOLEAN NOT NULL DEFAULT TRUE,
			conn_type TEXT NOT NULL
		);
		`,

		// rooms
		`
		CREATE TABLE IF NOT EXISTS rooms (
			id SERIAL PRIMARY KEY,
			owner_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// room_users
		`
		CREATE TABLE IF NOT EXISTS room_users (
			room_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (room_id, user_id),
			FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// room_admins
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
			id SERIAL PRIMARY KEY,
			user1_id INTEGER NOT NULL,
			user2_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE (user1_id, user2_id),
			FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// messages
		`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			room_id INTEGER,
			chat_id INTEGER,
			body TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
			FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
		);
		`,

		// connection_requests
		`
		CREATE TABLE IF NOT EXISTS connection_requests (
			id SERIAL PRIMARY KEY,
			from_user_id INTEGER NOT NULL,
			to_user_id INTEGER NOT NULL,
			message TEXT,
			status INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`,

		// friends
		`
		CREATE TABLE IF NOT EXISTS friends (
			user_id INTEGER NOT NULL,
			friend_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (user_id, friend_id),
			CHECK (user_id <> friend_id)
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

