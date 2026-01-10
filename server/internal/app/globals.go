package app

import (
	"sync"
)
var R GlobalRooms
var U GlobalUsers
var UsedIds []int
var (
	Sessions = make(map[string]int)
	SessionsMu sync.RWMutex
)
