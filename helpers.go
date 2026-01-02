package main

import (
	"math/rand/v2"
	"fmt"
)
func generateId() (int) {
	var cont bool
	var newId int
	for {
		cont = false
		newId = rand.IntN(1 << 32)
		for _, n := range usedIds {
			if n == newId {
				cont = true
			}
		}
		if !cont {
			break
		}
	}
	usedIds = append(usedIds, newId)
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
func messageExsists(room_id, mess_id, user_id int) (bool) {
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
func roomPublic(id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id && u.Type == TypePublic {
			return true
		}
	}
	return false
}
func roomExsists(id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id {
			return true
		}
	}
	return false
}
func roomExsistsToDelete(id, user_id int) (bool) {
	for _, u := range R.Rooms {
		if u.Id == id && u.UserId == user_id {
			return true
		}
	}
	return false
}
func userExsists(id int) (bool) {
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
func getUserById(user_id int) (User, error) {
	for _, u := range U.Users {
		if u.Id == user_id {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("Invalid user id")
}
