package core_models

type Kill struct {
	KillerID       string
	VictimID       string
	VictimTeam     string
	KillerTeam     string
	Time           int64
	KillerAircraft int
	VictimAircraft int
	Weapon         int
	ID             string
}

type Death struct {
	VictimID       string
	Time           int64
	VictimAircraft int
	ID             string
}

type LimitedUserData struct {
	ID         string   `json:"id"`
	PilotNames []string `gorm:"serializer:json" json:"pilotNames"`
	Kill       int      `json:"kills"`
	Deaths     int      `json:"deaths"`
	ELO        int      `json:"elo"`
	Rank       int      `json:"rank"`
}

type User struct {
	UID         string
	ID          string             `json:"id"`
	PilotNames  []string           `gorm:"serializer:json" json:"pilotNames"`
	LoginTimes  []int64            `gorm:"serializer:json" json:"loginTimes"`
	LogoutTimes []int64            `gorm:"serializer:json" json:"logoutTimes"`
	Kills       int                `json:"kills"`
	Deaths      int                `json:"deaths"`
	Spawns      []SpawnStruct      `gorm:"serializer:json" json:"spawns"`
	ELO         int                `json:"elo"`
	ELOHistory  []ELOHistoryStruct `gorm:"serializer:json" json:"eloHistory"`
	Rank        int                `json:"rank"`
}

type ELOHistoryStruct struct {
	ELO  int   `json:"elo"`
	Time int64 `json:"time"`
}
type SpawnStruct struct {
	AV42c   int `json:"0"`
	F26b    int `json:"1"`
	F45A    int `json:"2"`
	AH94    int `json:"3"`
	Invalid int `json:"4"`
}
