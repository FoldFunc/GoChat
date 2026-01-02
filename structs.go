package main

var usedIds []int
type Message struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	Body   string `json:"body"`
}
type Type string
const (
	TypePublic  Type = "public"
	TypePrivate Type = "private"
)
type Room struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
	Type     Type      `json:"type"`
	Admins   []User    `json:"admins"`
	Users    []User    `json:"users"`
}
type User struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Rooms []Room `json:"rooms"`
}
type GlobalMessages struct {
	Messages []Message `json:"messages"`
}
type GlobalUsers struct {
	Users []User `json:"users"`
}
type GlobalRooms struct {
	Rooms []Room `json:"rooms"`
}
type NewRoomReq struct {
	RoomId   int    `json:"room_id"`
	RoomName string `json:"room_name"`
	UserId   int    `json:"user_id"`
	RoomType bool   `json:"room_type"`
}
type NewUserReq struct {
	UserName string `json:"user_name"`
}
type SendMessreq struct {
	RoomId int    `json:"room_id"`
	UserId int    `json:"user_id"`
	Body   string `json:"body"`
}
type RemovemesReq struct {
	RoomId int    `json:"room_id"`
	UserId int    `json:"user_id"`
	MessId int    `json:"mess_id"`
}
type RemoveRoomReq struct {
	RoomId int    `json:"room_id"`
	UserId int    `json:"user_id"`
}
type AddToCloseRoomReq struct {
	RoomId  int `json:"room_id"`
	AdminId int `json:"admin_id"`
	UserId  int `json:"user_id"`
}
type AccesRoomReq struct {
	RoomId int `json:"room_id"`
	UserId int `json:"user_id"`
}
