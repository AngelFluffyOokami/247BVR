package dbengine

// KillData todo
type KillData struct {
	PID    string      `json:"pid"`
	Victim PlayerEvent `json:"victim"`
	Killer PlayerEvent `json:"killer"`
	Weapon
	ServerInfo ServerInfo `json:"serverInfo"`
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

// XYZ todo
type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Weapon todo
type Weapon struct {
	Weapon                  string `json:"weapon"`
	WeaponUUID              string `json:"weaponUuid"`
	PreviousDamagedByUserID string `json:"previousDamagedByUserId"`
	PreviousDamagedByWeapon string `json:"previousDamagedByWeapon"`
}

// ServerInfo todo
type ServerInfo struct {
	OnlineUsers []string `json:"onlineUsers"`
	TimeOfDay   string   `json:"timeOfDay"`
	MissionID   string   `json:"missionId"`
}
