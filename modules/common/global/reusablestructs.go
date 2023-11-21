package global

import "time"

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
type WsPlayerEvent struct {
	OwnerID   string   `json:"ownerId"`
	Occupants []string `json:"occupants"`
	Position  XYZ      `json:"position"`
	Velocity  XYZ      `json:"velocity"`
	Team      string   `json:"team"`
	Type      string   `json:"type"`
}

type WsServerInfo struct {
	OnlineUsers []string `json:"onlineUsers"`
	TimeOfDay   string   `json:"timeOfDay"`
	MissionID   string   `json:"missionId"`
}
type PlayerEvent struct {
	OwnerID   string   `json:"ownerId"`
	Occupants []string `json:"occupants"`
	Position  XYZ      `json:"position"`
	Velocity  XYZ      `json:"velocity"`
	Team      int      `json:"team"`
	Type      int      `json:"type"`
}
type ServerInfo struct {
	OnlineUsers []string `json:"onlineUsers"`
	TimeOfDay   int      `json:"timeOfDay"`
	MissionID   string   `json:"missionId"`
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

type Spawns struct {
	A42 int `json:"0"`
	F26 int `json:"1"`
	F45 int `json:"2"`
	A94 int `json:"3"`
	INV int `json:"4"`
	T55 int `json:"5"`
}

type EloHistory struct {
	Elo  float64 `json:"elo"`
	Time int64   `json:"time"`
}

type EndOfSeasonStats struct {
	Season    int       `json:"season"`
	Elo       float64   `json:"elo"`
	Rank      int       `json:"rank"`
	TeamKills int       `json:"teamKills"`
	History   time.Time `json:"history"`
}

// XYZ todo
type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
