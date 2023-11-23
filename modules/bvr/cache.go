package bvr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/angelfluffyookami/HSVRUSB/modules/common/global"
	"github.com/angelfluffyookami/HSVRUSB/modules/common/utils/logger"
)

var firstBootDone bool

var InitDone = make(chan bool)

var logging = logger.Log{}

var PauseCache = make(chan bool)

func InitCache() {

	go syncCache()

	tick := time.NewTicker(2 * time.Minute)

	for {
		select {
		case <-tick.C:
			go syncCache()
		case pause := <-PauseCache:
			switch pause {
			case true:
				tick.Stop()
			case false:
				tick.Reset(2 * time.Minute)
			}
		}
	}

}

func syncCache() {
	PauseCache <- true

	go syncKills(0)
	go syncUsers(0)
	go syncDeaths(0)
	go syncOnline(0)

}

var killSyncDone = make(chan bool)
var deathSyncDone = make(chan bool)
var userSyncDone = make(chan bool)

func syncKills(retryCount int64) {
	var hadError bool

	req, err := http.Get(global.Config.APIEndpoint + "/kills")

	if err != nil {
		if retryCount == 4 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /kills. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
		return
	}

	TimeMark := time.Now()

	err = json.NewDecoder(req.Body).Decode(&Cache.Kills.Kills)
	TimeSince := time.Since(TimeMark)
	fmt.Println("Unmarshal Exec Time: " + TimeSince.String())
	if err != nil {
		if retryCount == 4 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /kills. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
		return
	}

	Cache.Kills.Timestamp = time.Now()

	defer func() {
		if hadError {
			logging.Err().Message("/kills sanity check fail, retrying.").Add()
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncKills(retryCount + 1)
			tick.Stop()
		}
	}()

	if !hadError {
		killSyncDone <- true

	}
}

func syncUsers(retryCount int64) {
	var hadError bool

	req, err := http.Get(global.Config.APIEndpoint + "/users")

	if err != nil {
		if retryCount == 4 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /users. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
		return
	}

	err = json.NewDecoder(req.Body).Decode(&Cache.Users.Users)
	if err != nil {
		if retryCount == 4 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /users. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
		return
	}

	Cache.Users.Timestamp = time.Now()

	defer func() {
		if hadError {
			logging.Err().Message("/users sanity check fail, retrying.").Add()
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncUsers(retryCount + 1)
			tick.Stop()
		}
	}()

	if !hadError {
		<-killSyncDone
		userSyncDone <- true
	}
}

func syncDeaths(retryCount int64) {
	var hadError bool

	req, err := http.Get(global.Config.APIEndpoint + "/deaths")

	if err != nil {
		if retryCount == 4 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /deaths. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
		return
	}

	err = json.NewDecoder(req.Body).Decode(&Cache.Deaths.Deaths)
	if err != nil {
		if retryCount == 4 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /deaths. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
		return
	}

	Cache.Deaths.Timestamp = time.Now()

	defer func() {
		if hadError {
			logging.Err().Message("/deaths sanity check fail, retrying.").Add()
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncDeaths(retryCount + 1)
			tick.Stop()
		}
	}()

	if !hadError {
		<-userSyncDone
		deathSyncDone <- true
	}

}

func syncOnline(retryCount int64) {

	var hadError bool
	req, err := http.Get(global.Config.APIEndpoint + "/online")

	if err != nil {
		if retryCount == 4 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /online. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
		return
	}

	err = json.NewDecoder(req.Body).Decode(&Cache.Online.Online)
	if err != nil {
		if retryCount == 4 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /online. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
		return
	}

	Cache.Online.Timestamp = time.Now()

	defer func() {
		if hadError {
			logging.Err().Message("/online sanity check fail, retrying.").Add()
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncOnline(retryCount + 1)
			tick.Stop()
		}
	}()
	if !hadError {
		<-deathSyncDone
		PauseCache <- false
		if !firstBootDone {
			firstBootDone = true
			InitDone <- true
		}
	}
}
