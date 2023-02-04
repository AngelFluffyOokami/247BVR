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
	ID         string
	PilotNames []string
	Kill       int
	Deaths     int
	ELO        int
	Rank       int
}

type User struct {
	ID          string
	UID         string
	PilotNames  []string
	LoginTimes  []int64
	LogoutTimes []int64
	Kills       int
	Deaths      int
	Spawns      SpawnStruct `gorm:"serializer:json"`
	ELO         int
	ELOHistory  ELOHistoryStruct `gorm:"serializer:json"`
	Rank        int
}

type ELOHistoryStruct struct {
	ELO  int
	Time int64
}
type SpawnStruct struct {
	Aircraft int
	Spawn    int
}
