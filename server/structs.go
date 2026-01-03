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
type ConnReq struct {
	FromReqId int    `json:"from_req_id"`
	Message   string `json:"message"`
}
type Chat struct {
	User1    User      `json:"user_1"`
	User2    User      `json:"user_2"`
	Messages []Message `json:"messages"`
}
type User struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Rooms        []Room    `json:"rooms"`
	Chats        []Chat    `json:"chats"`
	ConnType     Type    	 `json:"conn_type"`
	ConnRequests []ConnReq `json:"conn_request"`
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
	RoomName string `json:"room_name"`
	UserId   int    `json:"user_id"`
	RoomType bool   `json:"room_type"`
}
type NewUserReq struct {
	UserName string `json:"user_name"`
	ConnType bool 	`json:"conn_type"`
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
type SendUserReq struct {
	UserId  int    `json:"user_id"`
	SendId  int    `json:"send_id"`
	Message string `json:"message"`
}
type ViewUserReq struct {
	UserId int `json:"user_id"`
}
type GetNameByIdReq struct {
	UserId   int `json:"user_id"`
	SearchId int `json"search_id"`
}
