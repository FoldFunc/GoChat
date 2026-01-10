package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/FoldFunc/GoChat/server/internal/app"
)

func UserExists(userID int) (bool, error) {
	var exists int
	err := DB.QueryRow(
		`SELECT 1 FROM users WHERE id = ?`,
		userID,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		log.Println("Error: ", err)
		return false, nil
	}
	if err != nil {
		log.Println("Error: ", err)
		return false, err
	}
	return true, nil
}
func GetUserPasswordByIdDB(userID int) (string, error) {
	query := `
		SELECT password
		FROM users
		WHERE id = ?;
	`

	var user string

	err := DB.QueryRow(query, userID).Scan(
		&user,
	)
	if err != nil {
		if err == sql.ErrNoRows {
		  log.Println("Error in quering the password: ", err)
			return "", errors.New("user not found")
		}
		log.Println("Error in quering the password: ", err)
		return "", err
	}

	return user, nil
}
func GetUserByIdDB(userID int) (app.UserData, error) {
	query := `
		SELECT id, name, logged_in, conn_type
		FROM users
		WHERE id = ?;
	`

	var user app.UserData

	err := DB.QueryRow(query, userID).Scan(
		&user.Id,
		&user.Name,
		&user.LoggedIn,
		&user.ConnType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return app.UserData{}, errors.New("user not found")
		}
		return app.UserData{}, err
	}

	return user, nil
}

func RoomExistsDB(roomID int) (bool, error) {
	var exists int
	err := DB.QueryRow(
		`SELECT 1 FROM rooms WHERE id = ? LIMIT 1;`,
		roomID,
	).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
func IsRoomPublicDB(roomID int) (bool, error) {
	var roomType string

	err := DB.QueryRow(
		`SELECT type FROM rooms WHERE id = ? LIMIT 1;`,
		roomID,
	).Scan(&roomType)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("room not found")
		}
		return false, err
	}
	if roomType == "private" {
		return false, nil
	}
	return true, nil
}
func GetRoomByIDDB(roomID int) (app.RoomData, error) {
	query := `
		SELECT id, owner_id, name, type
		FROM rooms
		WHERE id = ?
		LIMIT 1;
	`

	var room app.RoomData

	err := DB.QueryRow(query, roomID).Scan(
		&room.Id,
		&room.UserId, // owner_id
		&room.Name,
		&room.Type,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return app.RoomData{}, errors.New("room not found")
		}
		return app.RoomData{}, err
	}

	return room, nil
}
func IsUserInRoomDB(userID, roomID int) (bool, error) {
	var exists int
	err := DB.QueryRow(
		`SELECT 1 FROM room_users WHERE room_id = ? AND user_id = ? LIMIT 1;`,
		roomID, userID,
	).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
func IsUserAdminInRoomDB(userID, roomID int) (bool, error) {
	var exists int
	err := DB.QueryRow(
		`SELECT 1 FROM room_admins WHERE room_id = ? AND user_id = ? LIMIT 1;`,
		roomID, userID,
	).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("error: ", err)
			return false, nil
		}
		log.Println("error: ", err)
		return false, err
	}

	return true, nil
}
func MessageExistsDB(roomID, messageID, userID int) (bool, error) {
	var exists int
	err := DB.QueryRow(
		`SELECT 1 FROM messages WHERE id = ? AND room_id = ? AND user_id = ? LIMIT 1;`,
		messageID, roomID, userID,
	).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func IsUserPrivateDB(userID int) (bool, error) {
	var connType int

	err := DB.QueryRow(
		`SELECT conn_type FROM users WHERE id = ? LIMIT 1;`,
		userID,
	).Scan(&connType)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, errors.New("user not found")
		}
		return false, err
	}
	return connType == 0, nil
}

