package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	"github.com/angelfluffyookami/247BVR/modules/dbengine"
	"github.com/bwmarrin/discordgo"
)

// OnLogin todo
func OnLogin(r *discordgo.Ready) {

}

func AssertValue(timeOfDay string) int {

	var tod int

	switch timeOfDay {
	case "Morning":
		tod = 0
	case "Day":
		tod = 1

	case "Night":
		tod = 2
	default:
		tod = 3
	}

	return tod
}

func AssertWeapons(weapon string) int {
	if strings.Contains(weapon, "GUN") {
		return 0
	} else if (strings.Contains(weapon, "AIM-120D")) || (strings.Contains(weapon, "AIM-120")) {
		return 1
	} else if strings.Contains(weapon, "AIM-9") {
		return 2
	} else if strings.Contains(weapon, "AIM-7") {
		return 3
	} else if strings.Contains(weapon, "AIM-9+") {
		return 4
	} else if strings.Contains(weapon, "AIRS-T") {
		return 5
	} else if strings.Contains(weapon, "HARM") {
		return 6
	} else if strings.Contains(weapon, "SideARM") {
		return 6
	} else if strings.Contains(weapon, "AIM-9E") {
		return 8
	} else if strings.Contains(weapon, "CFIT") {
		return 9
	} else if strings.Contains(weapon, "COLLISION") {
		return 10
	} else {
		return 7
	}

}

func AssertTeam(team string) int {
	switch team {
	case "Allied":
		return 0
	case "Enemy":
		return 1
	default:
		return 2
	}
}
func AssertAircraft(aircraft string) int {
	switch aircraft {
	case "Vehicles/VTOL4":
		return 0
	case "Vehicles/FA-26B":
		return 1
	case "Vehicles/SEVTF":
		return 2
	case "Vehicles/AH-94":
		return 3
	case "Vehicles/T-55":
		return 4
	default:
		return 5
	}
}

// OnKillStream todo
func OnKillStream(kill global.WsKillEvent) {

	var killDb = global.KillEvent{
		ID:                      "",
		ID0:                     "",
		Time:                    0,
		Weapon:                  AssertWeapons(kill.Weapon),
		WeaponUUID:              kill.WeaponUUID,
		PreviousDamagedByUserID: kill.PreviousDamagedByUserID,
		PreviousDamagedByWeapon: kill.PreviousDamagedByWeapon,
		Killer: global.PlayerEvent{
			OwnerID:   kill.Killer.OwnerID,
			Occupants: kill.Killer.Occupants,
			Position:  kill.Killer.Position,
			Velocity:  kill.Killer.Velocity,
			Team:      AssertTeam(kill.Killer.Team),
			Type:      AssertAircraft(kill.Killer.Type),
		},
		Victim: global.PlayerEvent{
			OwnerID:   kill.Victim.OwnerID,
			Occupants: kill.Victim.Occupants,
			Position:  kill.Victim.Position,
			Velocity:  kill.Victim.Velocity,
			Team:      AssertTeam(kill.Victim.Team),
			Type:      AssertAircraft(kill.Victim.Type),
		},
		ServerInfo: global.ServerInfo{
			OnlineUsers: kill.ServerInfo.OnlineUsers,
			TimeOfDay:   AssertValue(kill.ServerInfo.TimeOfDay),
			MissionID:   kill.ServerInfo.MissionID,
		},
		Identified: false,
		Season:     kill.Season,
	}

	dbengine.DBv.WriteDB("kill", killDb, killDb.WeaponUUID)
}

func OnOnlineStream(online []global.WsOnlineData) {
	pid := time.Now().Unix()
	dbengine.DBv.WriteDB("online", online, fmt.Sprint(pid))
}

func OnSpawnStream(spawn global.WsSpawnData) {
	pid := time.Now().Unix()
	dbengine.DBv.WriteDB("spawn", spawn, fmt.Sprint(pid))
}

func OnLoginStream(login global.WsUserLogEvent) {
	pid := time.Now().Unix()
	dbengine.DBv.WriteDB("login", login, fmt.Sprint(pid))
}

func OnLogoutStream(logout global.WsUserLogEvent) {
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
	endpointKills := getKillsJson("http://hs.vtolvr.live/api/v1/public/kills")

	if endpointKills == nil {
		return
	}

	databaseString, err := dbengine.DBv.Db.ReadAll("kill")

	if err != nil {
		log.Panic(err)
	}

	var databaseKills global.Kills

	for _, v := range databaseString {
		var x global.KillEvent
		json.Unmarshal([]byte(v), &x)
		databaseKills = append(databaseKills, x)
	}

	updateKill(endpointKills, databaseKills)
	populateKills(endpointKills, databaseKills)

}

func populateKills(endpointMapKills map[string]global.KillEvent, databaseKills global.Kills) {

	databaseMapKills := make(map[string]global.KillEvent)
	// Populate Map
	for _, v := range databaseKills {

		databaseMapKills[v.WeaponUUID] = v

	}

	for _, v := range endpointMapKills {
		_, ok := databaseMapKills[v.WeaponUUID]
		if !ok {
			v.Identified = true
			go dbengine.DBv.WriteDB("kill", v, v.WeaponUUID)
		}
	}

}

func getKillsJson(url string) map[string]global.KillEvent {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	killMap := make(map[string]global.KillEvent)

	for dec.More() {
		var kill global.KillEvent

		if err = dec.Decode(&kill); err != nil {
			log.Panic(err)
		}
		killMap[kill.WeaponUUID] = kill
	}

	if err != nil {

		log.Fatal(err)
		return nil
	}
	return killMap
}

func updateKill(endpointMapKills map[string]global.KillEvent, databaseKills global.Kills) {

	unidentifiedKills := make(map[string]global.KillEvent)
	// Populate Map
	for _, v := range databaseKills {
		if !v.Identified {
			unidentifiedKills[v.WeaponUUID] = v
		}
	}

	// Iterate, identify, and overwrite unidentified objects

	for _, v := range unidentifiedKills {
		v.ID = endpointMapKills[v.WeaponUUID].ID
		v.ID0 = endpointMapKills[v.WeaponUUID].ID0
		v.Time = endpointMapKills[v.WeaponUUID].Time
		v.Identified = true
		dbengine.DBv.WriteDB("kill", v, v.WeaponUUID)
	}

}
