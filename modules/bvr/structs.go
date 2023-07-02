package bvr

type AircraftStruct struct {
	AV42c   int
	FA26b   int
	F45A    int
	AH94    int
	Invalid int
}

var AircraftConst = AircraftStruct{
	AV42c:   1,
	FA26b:   2,
	F45A:    3,
	AH94:    4,
	Invalid: 0,
}

type WeaponStruct struct {
	Gun     int
	AIM120  int
	AIM9    int
	AIM7    int
	AIM9X   int
	AIRST   int
	HARM    int
	AIM9E   int
	Invalid int
}

var WeaponConst = WeaponStruct{
	Gun:     1,
	AIM120:  2,
	AIM9:    3,
	AIM7:    4,
	AIM9X:   5,
	AIRST:   6,
	HARM:    7,
	AIM9E:   8,
	Invalid: 0,
}

type TeamStruct struct {
	Allied  int
	Enemy   int
	Invalid int
}

var TeamConst = TeamStruct{
	Allied:  1,
	Enemy:   2,
	Invalid: 0,
}

type OnlineStruct struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
	Team string `json:"team,omitempty"`
}

type KillStruct struct {
	ID string `json:"id,omitempty"`

	Killer       KillerStruct     `json:"killer,omitempty"`
	VictimStruct VictimStruct     `json:"victim,omitempty"`
	ServerInfo   ServerInfoStruct `json:"serverInfo,omitempty"`

	Time int64 `json:"time,omitempty"`
}

type ServerInfoStruct struct {
	MissionID   string   `json:"missionId,omitempty"`
	OnlineUsers []string `json:"onlineUsers,omitempty"`
	TimeOfDay   int      `json:"timeOfDay,omitempty"`
}

type KillerStruct struct {
	OwnerID      string   `json:"ownerId,omitempty"`
	Occupants    []string `json:"occupants,omitempty"`
	AircraftType int      `json:"type,omitempty"`
	Team         int      `json:"team,omitempty"`
}

type VictimStruct struct {
	OwnerID      string   `json:"ownerId,omitempty"`
	Occupants    []string `json:"occupants,omitempty"`
	AircraftType int      `json:"type,omitempty"`
	Team         int      `json:"team,omitempty"`
}
type DeathStruct struct {
	ID string `json:"id,omitempty"`

	Season     int              `json:"season,omitempty"`
	ServerInfo ServerInfoStruct `json:"serverInfo,omitempty"`
	Victim     VictimStruct     `json:"victim,omitempty"`

	Time int64 `json:"Time"`
}

type LimitedUserDataStruct struct {
	ID         string   `json:"id,omitempty"`
	PilotNames []string `json:"pilotNames,omitempty"`
	Kill       int      `json:"kills,omitempty"`
	Deaths     int      `json:"deaths,omitempty"`
	ELO        float64  `json:"elo,omitempty"`
	Rank       int      `json:"rank,omitempty"`
}

type UserStruct struct {
	UID         string
	ID          string             `json:"id,omitempty"`
	PilotNames  []string           `json:"pilotNames,omitempty"`
	LoginTimes  []int64            `json:"loginTimes,omitempty"`
	LogoutTimes []int64            `json:"logoutTimes,omitempty"`
	Kills       int                `json:"kills,omitempty"`
	Deaths      int                `json:"deaths,omitempty"`
	TeamKills   int                `json:"teamKills,omitempty"`
	Banned      bool               `json:"isBanned,omitempty"`
	DiscordID   string             `json:"discordId,omitempty"`
	Spawns      SpawnStruct        `json:"spawns,omitempty"`
	ELO         float64            `json:"elo,omitempty"`
	ELOHistory  []ELOHistoryStruct `json:"eloHistory,omitempty"`
	Rank        int                `json:"rank,omitempty"`
}

type ELOHistoryStruct struct {
	ELO  float64 `json:"elo,omitempty"`
	Time int64   `json:"time,omitempty"`
}

type SpawnStruct struct {
	AV42c   int `json:"0,omitempty"`
	F26b    int `json:"1,omitempty"`
	F45A    int `json:"2,omitempty"`
	AH94    int `json:"3,omitempty"`
	T55     int `json:"4,omitempty"`
	Invalid int `json:"5,omitempty"`
}

type EndOfSeasonStruct struct {
	Season    int     `json:"season,omitempty"`
	ELO       float64 `json:"elo,omitempty"`
	Rank      int     `json:"rank,omitempty"`
	TeamKills int     `json:"teamKills,omitempty"`
	History   string  `json:"history,omitempty"`
}
