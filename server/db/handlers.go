package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/FoldFunc/GoChat/server/app"
)
func CreateUser(user app.UserData, password string) error {
	query := `INSERT INTO users (id, name, password, conn_type) VALUES (?, ?, ?, ?);`
	_, err := DB.Exec(query, user.Id, user.Name, password, user.ConnType)
	if err != nil {
		log.Println("Error in creating user at db level: ", err)
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func CreateRoom(room app.RoomData) error {
	query := `INSERT INTO rooms (id, owner_id, name, type) VALUES (?, ?, ?, ?);`
	_, err := DB.Exec(query, room.Id, room.UserId, room.Name, room.Type)
	if err != nil {
		fmt.Printf("Error in insertion: %e", err)
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func InsertMessageRoom(message app.Message, room app.RoomData) error {
	query := `INSERT INTO messages (user_id, room_id, chat_id, body) VALUES (?, ?, NULL, ?);`
	_, err := DB.Exec(query, message.UserId, room.Id, message.Body)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func InsertUserCloseRoom(user app.UserData, room app.RoomData) error {
	query := `INSERT INTO room_users (room_id, user_id) VALUES (?, ?);`
	_, err := DB.Exec(query, room.Id, user.Id)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func RemoveMessage(user app.UserData, room app.RoomData, message int) error {
	query := `DELETE FROM messages WHERE id = ? AND room_id = ?;`
	_, err := DB.Exec(query, message, room.Id)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func RemoveRoom(room app.RoomData) error {
	query := `DELETE FROM rooms WHERE id = ?;`
	_, err := DB.Exec(query, room.Id)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func AddUserReq(conn app.ConnReq, toUser int) error {
	query := `INSERT INTO connection_requests (from_user_id, to_user_id, message) VALUES (?, ?, ?);`
	_, err := DB.Exec(query, conn.FromReqId, toUser, conn.Message)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
 	}
	return nil
}
func GetConnReq(user app.UserData) ([]app.ConnectionRequest, error){
	query := `
		SELECT id, from_user_id, to_user_id, message, status, created_at
		FROM connection_requests
		WHERE to_user_id = ? AND status = 0
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
		if err := rows.Scan(
			&r.ID,
			&r.FromUserID,
			&r.ToUserID,
			&r.Message,
			&r.Status,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, r)	
	}
	return requests, nil
}
func GetNameByIdDB(userID int) (string, error) {
	query := `SELECT name FROM users WHERE id = ?;`

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
		WHERE ru.user_id = ?;
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []app.RoomData

	for rows.Next() {
		var r app.RoomData
		if err := rows.Scan(
			&r.Id,
			&r.UserId,
			&r.Name,
			&r.Type,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
func QueryUserChatsDB(userID int) ([]app.ChatData, error) {
	query := `ha
		SELECT id, user1_id, user2_id
		FROM chats
		WHERE user1_id = ? OR user2_id = ?;
	`

	rows, err := DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []app.ChatData

	for rows.Next() {
		var chatID int
		var user1ID int
		var user2ID int

		if err := rows.Scan(&chatID, &user1ID, &user2ID); err != nil {
			return nil, err
		}

		chat := app.ChatData{
			Id: chatID,
			User1Id: user1ID,
			User2Id: user2ID,
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
func GetUserIdByNameDB(username string) (int, error) {
	query := `SELECT id FROM users WHERE name = ?;`

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
func GetUserNameByIdDB(userId int) (int, error) {
	query := `SELECT name FROM users WHERE id = ?;`

	var id int
	err := DB.QueryRow(query, userId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}

	return id, nil
}
func GetChatBetweenUsersDB(userA, userB int) (app.ChatData, error) {
	query := `
		SELECT id, user1_id, user2_id
		FROM chats
		WHERE
			(user1_id = ? AND user2_id = ?)
		   OR
			(user1_id = ? AND user2_id = ?)
		LIMIT 1;
	`

	var chatID int
	var user1ID int
	var user2ID int

	err := DB.QueryRow(query, userA, userB, userB, userA).
		Scan(&chatID, &user1ID, &user2ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return app.ChatData{}, errors.New("chat does not exist")
		}
		return app.ChatData{}, err
	}

	chat := app.ChatData{
		Id: chatID,
		User1Id: user1ID,
		User2Id: user2ID,
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
	query := `
		SELECT id
		FROM rooms
		WHERE name = ? AND TYPE = 'public'
		LIMIT 1;
	`

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
	log.Printf("userID: %d;roomName: %s\n", userID, roomName)
	query := `
		SELECT r.id, r.owner_id, r.name, r.type
		FROM rooms r
		JOIN room_users ru ON ru.room_id = r.id
		WHERE r.name = ? AND ru.user_id = ?
		LIMIT 1;
	`

	var room app.RoomData

	err := DB.QueryRow(query, roomName, userID).Scan(
		&room.Id,
		&room.UserId,
		&room.Name,
		&room.Type,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return app.RoomData{}, errors.New("room not found or access denied")
		}
		return app.RoomData{}, err
	}

	return room, nil
}
func SetUserLoggedInDB(userID int, loggedIn bool) error {
	var val int
	if loggedIn {
		val = 1
	} else {
		val = 0
	}

	result, err := DB.Exec(
		`UPDATE users SET logged_in = ? WHERE id = ?;`,
		val, userID,
	)
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
func AddToRoomUserDB(room app.RoomData, userId int) (error) {
	query := `INSERT INTO room_users (room_id, user_id) VALUES (?, ?);`
	_, err := DB.Exec(query, room.Id, userId)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}
	return nil
}
func AddUserAsAdminDB(room app.RoomData, userId int) (error) {
	query := `INSERT INTO room_admins (room_id, user_id) VALUES (?, ?);`
	_, err := DB.Exec(query, room.Id, userId)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}
	return nil
}
