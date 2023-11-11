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
	TrackingType string `json:"trackingType"`
	TrackingData any    `json:"trackingData"`
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
	ServerInfo ServerInfo `json:"serverInfo"`
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
	PilotName string `json:"pilotName,omitEmpty"`
}
