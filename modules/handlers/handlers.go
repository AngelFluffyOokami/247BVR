package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func Sync() {

	killSync()

}
func killSync() {
	msg := getJson("http://hs.vtolvr.live/api/v1/public/kills")

	if msg == "" {
		return
	}

	var kills []global.KillData

	json.Unmarshal([]byte(msg), &kills)

	compareKill(kills)

}

func getJson(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(resBody)
}

func compareKill(endpointKills []global.KillData) {
	var endpointMapKills map[string]global.KillData
	var unidentifiedKills map[string]global.KillData

	endpointMapKills = make(map[string]global.KillData)
	unidentifiedKills = make(map[string]global.KillData)

	databaseString, err := dbengine.DBv.Db.ReadAll("kill")

	if err != nil {
		return
	}

	var databaseKills []global.KillData

	for _, v := range databaseString {
		var x global.KillData
		json.Unmarshal([]byte(v), &x)
		databaseKills = append(databaseKills, x)
	}

	// Populate Map
	for _, v := range databaseKills {

		unidentifiedKills[v.WeaponUUID] = v

	}

	// Populate Map
	for _, v := range endpointKills {

		endpointMapKills[v.WeaponUUID] = v
	}

	// Iterate, identify, and overwrite unidentified objects
	for _, v := range unidentifiedKills {

		if endpointMapKills[v.WeaponUUID].Identified {

			dbengine.DBv.WriteDB("kill", endpointMapKills[v.WeaponUUID], endpointMapKills[v.WeaponUUID].WeaponUUID)

		}

	}

}
