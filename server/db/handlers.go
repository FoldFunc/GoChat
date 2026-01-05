package db

import (
	"fmt"
	"database/sql"
	"errors"
	"github.com/FoldFunc/GoChat/server/app"
)
func CreateUser(user app.User) error {
	query := `INSERT INTO user (id, name, password, conn_type) VALUES (?, ?, ?, ?);`

	_, err := DB.Exec(query, user.Id, user.Name, user.Password, user.ConnType)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func CreateRoom(room app.Room) error {
	query := `INSERT INTO rooms (id, owner_id, name) VALUES (?, ?, ?);`
	_, err := DB.Exec(query, room.Id, room.UserId, room.Name)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func InsertMessageRoom(message app.Message, room app.Room) error {
	query := `INSERT INTO messages (user_id, room_id, chat_id, body) VALUES (?, ?, NULL, ?);`
	_, err := DB.Exec(query, message.UserId, room.Id, message.Body)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func InsertUserCloseRoom(user app.User, room app.Room) error {
	query := `INSERT INTO room_users (room_id, user_id) VALUES (?, ?);`
	_, err := DB.Exec(query, room.Id, user.Id)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func RemoveMessage(user app.User, room app.Room, message int) error {
	query := `DELETE FROM messages WHERE id = ? AND room_id = ?;`
	_, err := DB.Exec(query, message, room.Id)
	if err != nil {
		return fmt.Errorf("Error in insertion: %e", err)
	}
	return nil
}
func RemoveRoom(room app.Room) error {
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
type ConnectionRequest struct {
	ID         int
	FromUserID int
	ToUserID   int
	Message    string
	Status     int
	CreatedAt  string
}
func GetConnReq(user app.User) ([]ConnectionRequest, error){
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
	var requests []ConnectionRequest
	for rows.Next() {
		var r ConnectionRequest
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

type RoomData struct {
	Id      int
	UserId  int
	Name    string
	Type    app.Type
}
func QuerUserRoomsDB(userID int) ([]RoomData, error) {
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

	var rooms []RoomData

	for rows.Next() {
		var r RoomData
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
type ChatData struct {
	Id      int
	User1Id int
	User2Id int
}
func QueryUserChatsDB(userID int) ([]ChatData, error) {
	query := `
		SELECT id, user1_id, user2_id
		FROM chats
		WHERE user1_id = ? OR user2_id = ?;
	`

	rows, err := DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []ChatData

	for rows.Next() {
		var chatID int
		var user1ID int
		var user2ID int

		if err := rows.Scan(&chatID, &user1ID, &user2ID); err != nil {
			return nil, err
		}

		chat := ChatData{
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
func GetChatBetweenUsersDB(userA, userB int) (ChatData, error) {
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
			return ChatData{}, errors.New("chat does not exist")
		}
		return ChatData{}, err
	}

	chat := ChatData{
		Id: chatID,
		User1Id: user1ID,
		User2Id: user2ID,
	}

	return chat, nil
}
func QuerySpecificUserChatDB(currentUserID int, otherUserName string) (ChatData, error) {
	otherUserID, err := GetUserIdByNameDB(otherUserName)
	if err != nil {
		return ChatData{}, err
	}

	return GetChatBetweenUsersDB(currentUserID, otherUserID)
}
func QueryUserRoomByNameDB(userID int, roomName string) (RoomData, error) {
	query := `
		SELECT r.id, r.owner_id, r.name, r.type
		FROM rooms r
		JOIN room_users ru ON ru.room_id = r.id
		WHERE r.name = ? AND ru.user_id = ?
		LIMIT 1;
	`

	var room RoomData

	err := DB.QueryRow(query, roomName, userID).Scan(
		&room.Id,
		&room.UserId,
		&room.Name,
		&room.Type,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return RoomData{}, errors.New("room not found or access denied")
		}
		return RoomData{}, err
	}

	return room, nil
}
