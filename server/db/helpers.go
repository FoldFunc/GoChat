package db

import (
	"database/sql"
	"errors"

	"github.com/FoldFunc/GoChat/server/internal/app"
)


func UserExists(userID int) (bool, error) {
	var exists int
	err := DB.QueryRow(`SELECT 1 FROM users WHERE id = $1;`, userID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetUserPasswordByIdDB(userID int) (string, error) {
	query := `SELECT password FROM users WHERE id = $1;`
	var password string
	err := DB.QueryRow(query, userID).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return password, nil
}

func GetUserByIdDB(userID int) (app.UserData, error) {
	query := `
	SELECT id, name, logged_in, conn_type
	FROM users
	WHERE id = $1;
	`
	var user app.UserData
	err := DB.QueryRow(query, userID).Scan(&user.Id, &user.Name, &user.LoggedIn, &user.ConnType)
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
	err := DB.QueryRow(`SELECT 1 FROM rooms WHERE id = $1 LIMIT 1;`, roomID).Scan(&exists)
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
	err := DB.QueryRow(`SELECT type FROM rooms WHERE id = $1 LIMIT 1;`, roomID).Scan(&roomType)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("room not found")
		}
		return false, err
	}
	return roomType != "private", nil
}

func GetRoomByIDDB(roomID int) (app.RoomData, error) {
	query := `
	SELECT id, owner_id, name, type
	FROM rooms
	WHERE id = $1
	LIMIT 1;
	`
	var room app.RoomData
	err := DB.QueryRow(query, roomID).Scan(&room.Id, &room.UserId, &room.Name, &room.Type)
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
	err := DB.QueryRow(`SELECT 1 FROM room_users WHERE room_id = $1 AND user_id = $2 LIMIT 1;`, roomID, userID).Scan(&exists)
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
	err := DB.QueryRow(`SELECT 1 FROM room_admins WHERE room_id = $1 AND user_id = $2 LIMIT 1;`, roomID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func MessageExistsDB(roomID, messageID, userID int) (bool, error) {
	var exists int
	err := DB.QueryRow(`SELECT 1 FROM messages WHERE id = $1 AND room_id = $2 AND user_id = $3 LIMIT 1;`, messageID, roomID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsUserPrivateDB(userID int) (bool, error) {
	var connType string
	err := DB.QueryRow(`SELECT conn_type FROM users WHERE id = $1 LIMIT 1;`, userID).Scan(&connType)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}
	return connType == "private", nil
}

