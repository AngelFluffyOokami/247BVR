package bvr

import "time"

var Cache CacheStruct

type CacheStruct struct {
	Kills  KillsStruct
	Users  UsersStruct
	Deaths DeathsStruct
	Online OnlinesStruct
}

type KillsStruct struct {
	Kills     []KillStruct
	Timestamp time.Time
}
type UsersStruct struct {
	Users     []UserStruct
	Timestamp time.Time
}
type DeathsStruct struct {
	Deaths    []DeathStruct
	Timestamp time.Time
}
type OnlinesStruct struct {
	Online    []OnlineStruct
	Timestamp time.Time
}
