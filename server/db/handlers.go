package db

import (
	"fmt"

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
