package global

type Data struct {
	Token         string
	AdminServer   string
	AdminChannel  string
	InfoChannel   string
	WarnChannel   string
	ErrChannel    string
	UpdateChannel string
	APIEndpoint   string
	Debugging     bool
	ActiveSession bool
}

// WsMessage defines the structure of a valid HS ELO Bot websocket base message.
type WsMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
	Pid  string `json:"pid"`
}

// Tracking todo
type Tracking struct {
	TrackingType string   `json:"trackingType"`
	TrackingData []string `json:"trackingData"`
}

// LookupData todo
type LookupData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Lookups todo
type Lookups struct{}

// Kills todo
func (Lookups) Kills() string {
	return "kill"
}

// Online todo
type Online struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Team string `json:"team"`
}

// PlayerEvent todo
type PlayerEvent struct {
	OwnerID   string   `json:"ownerId"`
	Occupants []string `json:"occupants"`
	Position  XYZ      `json:"position"`
	Velocity  XYZ      `json:"velocity"`
	Team      string   `json:"team"`
	Type      string   `json:"type"`
}

// KillData todo
type KillData struct {
	Victim PlayerEvent `json:"victim"`
	Killer PlayerEvent `json:"killer"`
	Weapon
	Identification
	Identified bool
	Season     int        `json:"season"`
	ServerInfo ServerInfo `json:"serverInfo"`
}

type SpawnPlace struct {
	Property1 int `json:"property1"`
	Property2 int `json:"property2"`
}

type EloHistory struct {
	Elo  int `json:"elo"`
	Time int `json:"time"`
}

type UserInfo struct {
	Id                      string        `json:"id"`
	PilotNames              []string      `json:"pilotNames"`
	LoginTimes              []int         `json:"loginTimes"`
	LogoutTimes             []int         `json:"logoutTimes"`
	Kills                   int           `json:"kills"`
	Deaths                  int           `json:"deaths"`
	Spawns                  SpawnPlace    `json:"spawns"`
	Elo                     int           `json:"elo"`
	EloHistory              []EloHistory  `json:"eloHistory"`
	Rank                    int           `json:"rank"`
	History                 []string      `json:"history"`
	DiscordId               string        `json:"discordId"`
	IsBanned                bool          `json:"isBanned"`
	TeamKills               int           `json:"teamKills"`
	IgnoreKillsAgainstUsers []string      `json:"ignoreKillsAgainstUsers"`
	EndOfSeasonStats        []SeasonStats `json:"endOfSeasonStats"`
}

type SeasonStats struct {
	Season    int    `json:"season"`
	Rank      int    `json:"rank"`
	Elo       int    `json:"elo"`
	TeamKills int    `json:"teamKills"`
	History   string `json:"history"`
}

type DeathData struct {
	Victim     PlayerEvent `json:"victim"`
	Killer     PlayerEvent `json:"killer"`
	ServerInfo ServerInfo  `json:"serverInfo"`
	Cause      string      `json:"cause"`
	Season     int         `json:"season"`
	Identification
	Identified bool
}

type Identification struct {
	Id   string `json:"id"`
	Time int    `json:"time"`
}

// Weapon todo
type Weapon struct {
	Weapon                  string `json:"weapon"`
	WeaponUUID              string `json:"weaponUuid"`
	PreviousDamagedByUserID string `json:"previousDamagedByUserId"`
	PreviousDamagedByWeapon string `json:"previousDamagedByWeapon"`
}

// XYZ todo
type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// ServerInfo todo
type ServerInfo struct {
	OnlineUsers []string `json:"onlineUsers"`
	TimeOfDay   string   `json:"timeOfDay"`
	MissionID   string   `json:"missionId"`
}

// SpawnData todo
type SpawnData struct {
	Player     PlayerEvent `json:"user"`
	ServerInfo ServerInfo  `json:"serverInfo"`
}

// OnlineData todo
type OnlineData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Team string `json:"team"`
}

// UserLogEvent todo
type UserLogEvent struct {
	UserID    string `json:"userId"`
	PilotName string `json:"pilotName"`
}
