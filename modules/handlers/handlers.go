package handlers

import (
	"fmt"
	"time"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	"github.com/angelfluffyookami/247BVR/modules/dbengine"
	"github.com/bwmarrin/discordgo"
)

// OnLogin todo
func OnLogin(r *discordgo.Ready) {

}

// OnKillStream todo
func OnKillStream(kill global.KillData) {
	dbengine.DBv.WriteDB("kill", kill, kill.WeaponUUID)
}

func OnOnlineStream(online []global.OnlineData) {
	pid := time.Now().Unix()
	dbengine.DBv.WriteDB("online", online, fmt.Sprint(pid))
}

func OnSpawnStream(spawn global.SpawnData) {
	pid := time.Now().Unix()
	dbengine.DBv.WriteDB("spawn", spawn, fmt.Sprint(pid))
}

func OnLoginStream(login global.UserLogEvent) {
	pid := time.Now().Unix()
	dbengine.DBv.WriteDB("login", login, fmt.Sprint(pid))
}

func OnLogoutStream(logout global.UserLogEvent) {
	pid := time.Now().Unix()
	dbengine.DBv.WriteDB("logout", logout, fmt.Sprint(pid))

}

func OnTrackingStream(tracking global.Tracking) {
	pid := time.Now().Unix()

	dbengine.DBv.WriteDB(tracking.TrackingType, tracking.TrackingData, fmt.Sprint(pid))
}
