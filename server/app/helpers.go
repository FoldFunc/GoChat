package app

import (
	"fmt"
	"math/rand/v2"

)
func GenerateId() (int) {
	var cont bool
	var newId int
	for {
		cont = false
		newId = rand.IntN(1 << 32)
		for _, n := range UsedIds {
			if n == newId {
				cont = true
			}
		}
		if !cont {
			break
		}
	}
	UsedIds = append(UsedIds, newId)
	return newId
}
func IsAdmin(user_id, room_id int) (bool) {
	for _, r := range R.Rooms {
		if r.Id == room_id {
			for _, a := range r.Admins {
				if a.Id == user_id {
					return true
				}
			}
		}
	}
	return false
}
func MessageExsists(room_id, mess_id, user_id int) (bool) {
	for _, r := range R.Rooms {
		if r.Id == room_id {
			for _, m := range r.Messages {
				if m.Id == mess_id && m.UserId == user_id {
					return true
				}
			}
		}
	}
	return false
}
func RoomPublic(id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id && u.Type == TypePublic {
			return true
		}
	}
	return false
}
func RoomExsists(id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id {
			return true
		}
	}
	return false
}
func RoomExsistsToDelete(id, user_id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id && u.UserId == user_id {
			return true
		}
	}
	return false
}
func UserExsists(id int) (bool) {
	for _, u := range U.Users {
		if u.Id == id {
			return true
		}
	}
	return false
}
func UserInRoom(user_id, room_id int) (bool) {
	for _, u := range U.Users {
		if u.Id == user_id {
			for _, r := range u.Rooms {
				if r.Id == room_id {
					return true
				}
			}
		}
	}
	return false
}
func GetRoomById(room_id int) (*Room , error) {
	for _, r := range R.Rooms {
		if r.Id == room_id{
			return &r, nil
		}
	}
	return &Room{}, fmt.Errorf("Invalid room id")
}
func GetUserByName(user_name string) (*User, error) {
	for _, u := range U.Users {
		if u.Name == user_name {
			return &u, nil
		}
	}
	return &User{}, fmt.Errorf("Invalid user name")
}
func GetUserById(user_id int) (*User, error) {
	for _, u := range U.Users {
		if u.Id == user_id {
			return &u, nil
		}
	}
	return &User{}, fmt.Errorf("Invalid user id")
}
func UserPrivate(user_id int) (bool) {
	for _, u := range U.Users {
		if u.Id == user_id && u.ConnType == TypePrivate {
			return true
		}
	}
	return false
}
