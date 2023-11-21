package global

import "time"

type Kills []KillEvent

type KillEvent struct {
	ID                      string      `json:"_id"`
	ID0                     string      `json:"id"`
	Time                    int64       `json:"time"`
	Weapon                  int         `json:"weapon"`
	WeaponUUID              string      `json:"weaponUuid"`
	PreviousDamagedByUserID string      `json:"previousDamagedByUserId"`
	PreviousDamagedByWeapon int         `json:"previousDamagedByWeapon"`
	Killer                  PlayerEvent `json:"killer"`
	Identified              bool
	Victim                  PlayerEvent `json:"victim"`
	ServerInfo              ServerInfo  `json:"serverInfo"`
	Season                  int         `json:"season"`
}

type Deaths []DeathEvent

type DeathEvent struct {
	ID                      string      `json:"_id"`
	ID0                     string      `json:"id"`
	Time                    int64       `json:"time"`
	Weapon                  int         `json:"weapon"`
	WeaponUUID              string      `json:"weaponUuid"`
	PreviousDamagedByUserID string      `json:"previousDamagedByUserId"`
	PreviousDamagedByWeapon int         `json:"previousDamagedByWeapon"`
	Killer                  PlayerEvent `json:"killer"`
	Victim                  PlayerEvent `json:"victim"`
	ServerInfo              ServerInfo  `json:"serverInfo"`
	Season                  int         `json:"season"`
}

type Users []User

type User struct {
	ID               string             `json:"_id"`
	ID0              string             `json:"id"`
	PilotNames       []string           `json:"pilotNames"`
	LoginTimes       []int64            `json:"loginTimes"`
	LogoutTimes      []int64            `json:"logoutTimes"`
	Kills            int                `json:"kills"`
	Deaths           int                `json:"deaths"`
	Spawns           Spawns             `json:"spawns"`
	Elo              float64            `json:"elo"`
	EloHistory       []EloHistory       `json:"eloHistory"`
	DiscordID        string             `json:"discordId"`
	TeamKills        int                `json:"teamKills"`
	EndOfSeasonStats []EndOfSeasonStats `json:"endOfSeasonStats"`
	Rank             int                `json:"rank"`
	History          []time.Time        `json:"history"`
}

type UserLimited struct {
	ID         string   `json:"id"`
	PilotNames []string `json:"pilotNames"`
	Kills      int      `json:"kills"`
	Deaths     int      `json:"deaths"`
	Elo        float64  `json:"elo"`
	Rank       int      `json:"rank"`
	DiscordID  string   `json:"discordId"`
	TeamKills  int      `json:"teamKills"`
}
