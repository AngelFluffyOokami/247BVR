package global

type WsKillEvent struct {
	Victim                  WsPlayerEvent `json:"victim"`
	Killer                  WsPlayerEvent `json:"killer"`
	Weapon                  string        `json:"weapon"`
	WeaponUUID              string        `json:"weaponUuid"`
	PreviousDamagedByUserID string        `json:"previousDamagedByUserId"`
	PreviousDamagedByWeapon int           `json:"previousDamagedByWeapon"`
	Season                  int           `json:"season"`
	ServerInfo              WsServerInfo  `json:"serverInfo"`
}

type WsDeathData struct {
	Victim     WsPlayerEvent `json:"victim"`
	Killer     WsPlayerEvent `json:"killer"`
	ServerInfo WsServerInfo  `json:"serverInfo"`
	Cause      string        `json:"cause"`
	Season     int           `json:"season"`
}

// SpawnData todo
type WsSpawnData struct {
	Player     WsPlayerEvent `json:"user"`
	ServerInfo WsServerInfo  `json:"serverInfo"`
}

// OnlineData todo
type WsOnlineData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Team string `json:"team"`
}

// UserLogEvent todo
type WsUserLogEvent struct {
	UserID    string `json:"userId"`
	PilotName string `json:"pilotName"`
}

// Online todo
type WsOnline struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Team string `json:"team"`
}
