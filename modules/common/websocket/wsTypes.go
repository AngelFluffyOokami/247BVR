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
	ClosedChan         chan bool
	URI                string
	WsOnline           bool
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
