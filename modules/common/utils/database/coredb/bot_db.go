package core_models

type Cinnamon struct {
	BotID           string `gorm:"primaryKey"`
	Uptime          int64
	CreationDate    int64
	UpSince         int64
	TotalUptime     int64
	TotalDowntime   int64
	DowntimePercent float64
	PastUptime      []PastUptime `gorm:"serializer:json"`
}

type PastUptime struct {
	Uptime  int64
	UpSince int64
}
