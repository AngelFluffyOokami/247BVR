package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
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

	go userSync()

	go killSync()
	<-usersync
	<-killsync

	fmt.Println("Sync Done, continuing startup.")
}

var usersync = make(chan bool)
var killsync = make(chan bool)

func userSync() {

	limitedMap := usersJsonGet("https://hs.vtolvr.live/api/v1/public/users")
	empty := false
	userString, err := dbengine.DBv.Db.ReadAll("users")

	if err != nil {
		fmt.Println("err encountered, assuming empty")
		empty = true
	}
	var users global.Users

	if !empty {

		for _, v := range userString {
			var x global.User
			json.Unmarshal([]byte(v), &x)
			users = append(users, x)
		}
	}
	syncUser(limitedMap, users)
	usersync <- true
}

func syncUser(limitedUserMap map[string]global.UserLimited, users global.Users) {
	for _, v := range users {
		_, ok := limitedUserMap[v.ID0]
		if ok {
			delete(limitedUserMap, v.ID0)
		}
	}

	for _, v := range limitedUserMap {
		newUser := global.User{
			ID0:        v.ID0,
			PilotNames: v.PilotNames,
			Kills:      v.Kills,
			Deaths:     v.Deaths,
			Elo:        v.Elo,
			Rank:       v.Rank,
			DiscordID:  v.DiscordID,
			TeamKills:  v.TeamKills,
			Identified: false,
		}
		dbengine.DBv.WriteDB("users", newUser, newUser.ID0)
	}
	identifyUser()
}

func identifyUser() {

	dbUsers, _ := dbengine.DBv.Db.ReadAll("users")

	unidentifiedUsers := make(map[string]global.User)

	for _, v := range dbUsers {
		var user global.User
		json.Unmarshal([]byte(v), &user)
		if !user.Identified {
			unidentifiedUsers[user.ID0] = user

			if len(user.PilotNames) >= 1 {
				fmt.Println("Unidentified user: " + user.ID0 + " " + user.PilotNames[0])
			} else {
				fmt.Println("Unidentified user: " + user.ID0 + " " + user.ID)
			}

		}
	}

	for _, v := range unidentifiedUsers {
		threadSync.Add(1)
		go idThread(v)
		if currentThreads >= 64 {
			fmt.Println("waiting on threads")
			threadSync.Wait()
		}
		currentThreads += 1
		fmt.Println("current id threads: " + fmt.Sprint(currentThreads))

	}

	fmt.Println("waiting on last threads")
	threadSync.Wait()

}

var currentThreads = 0
var threadSync sync.WaitGroup

func idThread(unidentifiedUser global.User) {
	if len(unidentifiedUser.PilotNames) >= 1 {
		fmt.Println("Get Request: " + unidentifiedUser.ID0 + " " + unidentifiedUser.PilotNames[0])
	} else {
		fmt.Println("Get Request: " + unidentifiedUser.ID0 + " " + unidentifiedUser.ID)
	}
	newuser, ok := userJsonGet(unidentifiedUser.ID0, "https://hs.vtolvr.live/api/v1/public/users")
	if ok {
		go dbengine.DBv.WriteDB("users", newuser, newuser.ID0)

		if len(newuser.PilotNames) >= 1 {
			fmt.Println("Unidentified user: " + newuser.ID0 + " " + newuser.PilotNames[0] + "identified")
		} else {
			fmt.Println("Unidentified user: " + newuser.ID0 + " " + newuser.ID + "Identified")
		}
	} else {
		fmt.Println("BLEHHHH")
	}

	currentThreads -= 1
	defer threadSync.Done()

}
func userJsonGet(ID string, url string) (global.User, bool) {

	resp, err := http.Get(url + "/" + ID)
	if err != nil {
		return global.User{}, false
	}

	defer resp.Body.Close()

	msg, err := io.ReadAll(resp.Body)

	if err != nil {

		return global.User{}, false
	}

	var user global.User
	json.Unmarshal(msg, &user)
	if len(user.PilotNames) >= 1 {
		fmt.Println("Unmarshalled user: " + user.ID0 + " " + user.PilotNames[0])
	} else {
		fmt.Println("Unmarshalled user: " + user.ID0 + " " + user.ID)
	}
	return user, true

}

func usersJsonGet(url string) map[string]global.UserLimited {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	limitedMap := make(map[string]global.UserLimited)
	var limited []global.UserLimited

	json.Unmarshal(msg, &limited)
	for _, v := range limited {
		limitedMap[v.ID0] = v
	}

	if err != nil {

		log.Fatal(err)
		return nil
	}
	return limitedMap
}

func killSync() {
	endpointKills := getKillsJson("http://hs.vtolvr.live/api/v1/public/kills")

	if endpointKills == nil {
		return
	}

	databaseString, err := dbengine.DBv.Db.ReadAll("kill")

	empty := false
	if err != nil {
		fmt.Println("err returned, assuming db is empty")
		empty = true
	}

	var databaseKills global.Kills

	if !empty {

		for _, v := range databaseString {
			var x global.KillEvent
			json.Unmarshal([]byte(v), &x)
			databaseKills = append(databaseKills, x)
		}
	}

	populateKills(endpointKills, databaseKills)

	if !empty {
		updateKill(endpointKills, databaseKills)
	}
	killsync <- true
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
			fmt.Println("Populated kill: " + v.WeaponUUID + " " + fmt.Sprint(v.Weapon))
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

		if kill.WeaponUUID == "" {
			kill.WeaponUUID = kill.ID0
		}

		fmt.Println("Decoding kill JSON: " + kill.WeaponUUID + " " + fmt.Sprint(kill.Weapon))
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
		go dbengine.DBv.WriteDB("kill", v, v.WeaponUUID)

		fmt.Println("Identified Kill: " + v.WeaponUUID + " " + fmt.Sprint(v.Weapon))
	}

}
