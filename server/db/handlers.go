package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/FoldFunc/GoChat/server/internal/app"
)


func CreateUser(user app.UserData, password string) error {
	query := `INSERT INTO users (id, name, password, conn_type) VALUES ($1, $2, $3, $4);`
	_, err := DB.Exec(query, user.Id, user.Name, password, user.ConnType)
	if err != nil {
		log.Println("Error creating user at db level:", err)
		return fmt.Errorf("error in insertion: %w", err)
	}
	return nil
}

func CreateRoom(room app.RoomData) error {
	query := `INSERT INTO rooms (id, owner_id, name, type) VALUES ($1, $2, $3, $4);`
	_, err := DB.Exec(query, room.Id, room.UserId, room.Name, room.Type)
	if err != nil {
		return fmt.Errorf("error in insertion: %w", err)
	}
	return nil
}

func InsertMessageRoom(message app.Message, room app.RoomData) error {
	query := `INSERT INTO messages (user_id, room_id, chat_id, body) VALUES ($1, $2, NULL, $3);`
	_, err := DB.Exec(query, message.UserId, room.Id, message.Body)
	if err != nil {
		return fmt.Errorf("error in insertion: %w", err)
	}
	return nil
}

func InsertUserCloseRoom(user app.UserData, room app.RoomData) error {
	query := `INSERT INTO room_users (room_id, user_id) VALUES ($1, $2);`
	_, err := DB.Exec(query, room.Id, user.Id)
	if err != nil {
		return fmt.Errorf("error in insertion: %w", err)
	}
	return nil
}

func RemoveMessage(user app.UserData, room app.RoomData, message int) error {
	query := `DELETE FROM messages WHERE id = $1 AND room_id = $2;`
	_, err := DB.Exec(query, message, room.Id)
	if err != nil {
		return fmt.Errorf("error in deletion: %w", err)
	}
	return nil
}

func RemoveRoom(room app.RoomData) error {
	query := `DELETE FROM rooms WHERE id = $1;`
	_, err := DB.Exec(query, room.Id)
	if err != nil {
		return fmt.Errorf("error in deletion: %w", err)
	}
	return nil
}

func AddUserReq(conn app.ConnReq, toUser int) error {
	query := `INSERT INTO connection_requests (from_user_id, to_user_id, message) VALUES ($1, $2, $3);`
	_, err := DB.Exec(query, conn.FromReqId, toUser, conn.Message)
	if err != nil {
		return fmt.Errorf("error in insertion: %w", err)
	}
	return nil
}

func GetConnReq(user app.UserData) ([]app.ConnectionRequest, error) {
	query := `
	SELECT id, from_user_id, to_user_id, message, status, created_at
	FROM connection_requests
	WHERE to_user_id = $1 AND status = 0
	ORDER BY created_at DESC;
	`
	rows, err := DB.Query(query, user.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []app.ConnectionRequest
	for rows.Next() {
		var r app.ConnectionRequest
		if err := rows.Scan(&r.ID, &r.FromUserID, &r.ToUserID, &r.Message, &r.Status, &r.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}
	return requests, nil
}

func GetConnReqFrom(user app.UserData, userFrom app.UserData) ([]app.ConnectionRequest, error) {
	query := `
	SELECT id, from_user_id, to_user_id, message, status, created_at
	FROM connection_requests
	WHERE to_user_id = $1 AND status = 0 AND from_user_id = $2
	ORDER BY created_at DESC;
	`
	rows, err := DB.Query(query, user.Id, userFrom.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []app.ConnectionRequest
	for rows.Next() {
		var r app.ConnectionRequest
		if err := rows.Scan(&r.ID, &r.FromUserID, &r.ToUserID, &r.Message, &r.Status, &r.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}
	return requests, nil
}

func GetNameByIdDB(userID int) (string, error) {
	query := `SELECT name FROM users WHERE id = $1;`
	var name string
	err := DB.QueryRow(query, userID).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return name, nil
}

func QuerUserRoomsDB(userID int) ([]app.RoomData, error) {
	query := `
	SELECT r.id, r.owner_id, r.name, r.type
	FROM rooms r
	JOIN room_users ru ON ru.room_id = r.id
	WHERE ru.user_id = $1;
	`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []app.RoomData
	for rows.Next() {
		var r app.RoomData
		if err := rows.Scan(&r.Id, &r.UserId, &r.Name, &r.Type); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

func QueryUserChatsDB(userID int) ([]app.ChatData, error) {
	query := `
	SELECT id, user1_id, user2_id
	FROM chats
	WHERE user1_id = $1 OR user2_id = $2;
	`
	rows, err := DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []app.ChatData
	for rows.Next() {
		var chat app.ChatData
		if err := rows.Scan(&chat.Id, &chat.User1Id, &chat.User2Id); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func GetUserIdByNameDB(username string) (int, error) {
	query := `SELECT id FROM users WHERE name = $1;`
	var id int
	err := DB.QueryRow(query, username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return id, nil
}

func GetUserNameByIdDB(userID int) (string, error) {
	query := `SELECT name FROM users WHERE id = $1;`
	var name string
	err := DB.QueryRow(query, userID).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return name, nil
}

func GetChatBetweenUsersDB(userA, userB int) (app.ChatData, error) {
	query := `
	SELECT id, user1_id, user2_id
	FROM chats
	WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1)
	LIMIT 1;
	`
	var chat app.ChatData
	err := DB.QueryRow(query, userA, userB).Scan(&chat.Id, &chat.User1Id, &chat.User2Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.ChatData{}, errors.New("chat does not exist")
		}
		return app.ChatData{}, err
	}
	return chat, nil
}

func QuerySpecificUserChatDB(currentUserID int, otherUserName string) (app.ChatData, error) {
	otherUserID, err := GetUserIdByNameDB(otherUserName)
	if err != nil {
		return app.ChatData{}, err
	}
	return GetChatBetweenUsersDB(currentUserID, otherUserID)
}

func QueryPublicRoomByNameDB(roomName string) (int, error) {
	query := `SELECT id FROM rooms WHERE name = $1 AND type = 'public' LIMIT 1;`
	var room int
	err := DB.QueryRow(query, roomName).Scan(&room)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("room not found or access denied")
		}
		return -1, err
	}
	return room, nil
}

func QueryUserRoomByNameDB(userID int, roomName string) (app.RoomData, error) {
	query := `
	SELECT r.id, r.owner_id, r.name, r.type
	FROM rooms r
	JOIN room_users ru ON ru.room_id = r.id
	WHERE r.name = $1 AND ru.user_id = $2
	LIMIT 1;
	`
	var room app.RoomData
	err := DB.QueryRow(query, roomName, userID).Scan(&room.Id, &room.UserId, &room.Name, &room.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.RoomData{}, errors.New("room not found or access denied")
		}
		return app.RoomData{}, err
	}
	return room, nil
}

func SetUserLoggedInDB(userID int, loggedIn bool) error {
	query := `UPDATE users SET logged_in = $1 WHERE id = $2;`
	result, err := DB.Exec(query, loggedIn, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func AddToRoomUserDB(room app.RoomData, userId int) error {
	query := `INSERT INTO room_users (room_id, user_id) VALUES ($1, $2);`
	_, err := DB.Exec(query, room.Id, userId)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	return nil
}

func AddUserAsAdminDB(room app.RoomData, userId int) error {
	query := `INSERT INTO room_admins (room_id, user_id) VALUES ($1, $2);`
	_, err := DB.Exec(query, room.Id, userId)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}
	return nil
}

func GetRoomIDByNameDB(roomName string) (int, error) {
	query := `SELECT id FROM rooms WHERE name = $1;`
	var id int
	err := DB.QueryRow(query, roomName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("room not found")
		}
		return 0, err
	}
	return id, nil
}

func GetMessageIDByNameDB(roomID int, messageBody string, userID int) (int, error) {
	query := `
	SELECT id
	FROM messages
	WHERE room_id = $1 AND user_id = $2 AND body = $3
	ORDER BY id DESC
	LIMIT 1;
	`
	var messageID int
	err := DB.QueryRow(query, roomID, userID, messageBody).Scan(&messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("message not found")
		}
		return 0, err
	}
	return messageID, nil
}

func RequestExists(requestID, userId int) (bool, error) {
	query := `SELECT 1 FROM requests WHERE id = $1 AND to_user_id = $2 LIMIT 1;`
	var exists int
	err := DB.QueryRow(query, requestID, userId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func RequestAccept(userID int, requestID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var fromUserID int
	err = tx.QueryRow(`SELECT from_user_id FROM requests WHERE id = $1 AND to_user_id = $2 LIMIT 1;`, requestID, userID).Scan(&fromUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("request not found or not authorized")
		}
		return err
	}

	_, err = tx.Exec(`INSERT INTO friends (user_id, friend_id) VALUES ($1, $2);`, userID, fromUserID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM requests WHERE id = $1;`, requestID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
