package app

type ConnectionRequest struct {
	ID         int
	FromUserID int
	ToUserID   int
	Message    string
	Status     int
	CreatedAt  string
}
type RoomData struct {
	Id      int
	UserId  int
	Name    string
	Type    Type
}
type ChatData struct {
	Id      int
	User1Id int
	User2Id int
}
type UserData struct {
	Id       int
	Name     string
	LoggedIn bool
	ConnType bool
}

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
	Id       int       			`json:"id"`
	UserId   int       			`json:"user_id"`
	Name     string    			`json:"name"`
	Messages []Message 			`json:"messages"`
	Type     Type      			`json:"type"`
	Admins   []UserData `json:"admins"`
	Users    []UserData `json:"users"`
}
type ConnReq struct {
	FromReqId int    `json:"from_req_id"`
	Message   string `json:"message"`
}
type Chat struct {
	User1    UserData  `json:"user_1"`
	User2    UserData  `json:"user_2"`
	Messages []Message `json:"messages"`
}
type User struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	LoggedIn     bool      `json:"logged_in"`
	Rooms        []RoomData`json:"rooms"`
	Chats        []ChatData`json:"chats"`
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
	RoomType bool   `json:"room_type"`
}
type NewUserReq struct {
	UserName string `json:"user_name"`
	ConnType bool 	`json:"conn_type"`
	Password string `json:"password"`
}
type SendMessreq struct {
	RoomId int    `json:"room_id"`
	Body   string `json:"body"`
}
type RemovemesReq struct {
	RoomId int    `json:"room_id"`
	MessId int    `json:"mess_id"`
}
type RemoveRoomReq struct {
	RoomId int    `json:"room_id"`
}
type AddToCloseRoomReq struct {
	RoomId  int `json:"room_id"`
	UserId  int `json:"user_id"`
}
type AccesRoomReq struct {
	RoomId int `json:"room_id"`
}
type SendUserReq struct {
	SendId  int    `json:"send_id"`
	Message string `json:"message"`
}
type GetNameByIdReq struct {
	SearchId int `json:"search_id"`
}
type GetIdByNameReq struct {
	SearchName string `json:"search_name"`
}
type LoginReq struct {
	UserName string `json:"user_name"`
	UserPassword string `json:"user_password"`
}
type QueryUserChatReq struct {
	ChatWithName string `json:"chat_with_name"`
}
type QueryUserRoomReq struct {
	RoomName string `json:"room_name"`
}
type QueryMessageFromRoomReq struct {
	RoomId      int    `json:"room_id"`
	MessageBody string `json:"message_body"`
}
