package app

import (
	"math/rand/v2"

)
func GenerateId() (int) {
	var cont bool
	var newId int
	for {
		cont = false
		newId = rand.IntN(1 << 8)
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
