package wshandler

import (
	"context"
	"sync"

	"nhooyr.io/websocket"
)

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

// MessageType contains constants of valid message types.
type MessageType struct {
}

// Subscribe sets message type to subscribe
func (MessageType) Subscribe() string {
	return "subscribe"
}

// SubscribeData contains constants of valid subscription types.
type SubscribeData struct {
}

// Login subscribes to user_login
func (SubscribeData) Login() string {
	return "user_login"
}

// Logout subscribes to user_logout
func (SubscribeData) Logout() string {
	return "user_logout"
}

// Kill subscribes to kill
func (SubscribeData) Kill() string {
	return "kill"
}

// Death subscribes to death
func (SubscribeData) Death() string {
	return "death"
}

// Spawn subscribes to spawn
func (SubscribeData) Spawn() string {
	return "spawn"
}

// Tracking subscribes to tracking
func (SubscribeData) Tracking() string {
	return "tracking"
}

// Online subscribes to online
func (SubscribeData) Online() string {
	return "online"
}

// Daemon subscribes to daemon_report
func (SubscribeData) Daemon() string {
	return "daemon_report"
}

// Missile subscribes to missile_launch_params
func (SubscribeData) Missile() string {
	return "missile_launch_params"
}

// All subscribes to all
func (SubscribeData) All() []string {
	return []string{
		"missile_launch_params",
		"daemon_report",
		"online",
		"tracking",
		"spawn",
		"death",
		"kill",
		"user_logout",
		"user_login",
	}
}

// LookupType todo
func (*Ws) LookupType() string {
	return "lookup"
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

// Ws todo
type Ws struct {
	Con                context.Context
	C                  *websocket.Conn
	Subscriptions      SubscribeData
	Lookups            Lookups
	URI                string
	CloseWS            chan bool
	wg                 sync.WaitGroup
	writeWg            sync.Mutex
	last               string
	TypesFound         []string
	TrackingTypesFound []string
	GoodUnmarshals     []string
	BadUnmarshals      []string
	SecondPerAttempt   int
}

// PongMessage todo
type PongMessage struct {
	Type string `json:"type"`
	PID  string `json:"pid"`
}

// SubscribeType todo
func (*Ws) subscribeType() string {
	return "subscribe"
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
