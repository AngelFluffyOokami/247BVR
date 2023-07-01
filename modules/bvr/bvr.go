package bvr

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	"github.com/angelfluffyookami/247BVR/modules/common/utils/logger"
)

var logging = logger.Log{}

var PauseCache = make(chan bool)

func InitCache() {

	tick := time.NewTicker(2 * time.Minute)

	select {
	case <-tick.C:
		syncCache()
	case pause := <-PauseCache:
		for pause {
			pause = <-PauseCache
		}
	}

}

func syncCache() {

	go syncKills(0)
	go syncUsers(0)
	go syncDeaths(0)
	go syncOnline(0)

}

func syncKills(retryCount int64) {
	var hadError bool

	req, err := http.Get(global.Config.APIEndpoint + "/kills")

	if err != nil {
		if retryCount == 5 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /kills. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
	}

	err = json.NewDecoder(req.Body).Decode(&Cache.Kills.Kills)
	if err != nil {
		if retryCount == 5 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /kills. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
	}

	Cache.Kills.Timestamp = time.Now()

	defer func() {
		if hadError {
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncKills(retryCount + 1)
			tick.Stop()
		}
	}()

}

func syncUsers(retryCount int64) {
	var hadError bool

	req, err := http.Get(global.Config.APIEndpoint + "/users")

	if err != nil {
		if retryCount == 5 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /users. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
	}

	err = json.NewDecoder(req.Body).Decode(&Cache.Users.Users)
	if err != nil {
		if retryCount == 5 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /users. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
	}

	Cache.Users.Timestamp = time.Now()

	defer func() {
		if hadError {
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncKills(retryCount + 1)
			tick.Stop()
		}
	}()
}

func syncDeaths(retryCount int64) {
	var hadError bool

	req, err := http.Get(global.Config.APIEndpoint + "/deaths")

	if err != nil {
		if retryCount == 5 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /deaths. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
	}

	err = json.NewDecoder(req.Body).Decode(&Cache.Deaths.Deaths)
	if err != nil {
		if retryCount == 5 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /deaths. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
	}

	Cache.Deaths.Timestamp = time.Now()

	defer func() {
		if hadError {
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncKills(retryCount + 1)
			tick.Stop()
		}
	}()
}

func syncOnline(retryCount int64) {

	var hadError bool
	req, err := http.Get(global.Config.APIEndpoint + "/online")

	if err != nil {
		if retryCount == 5 {
			PauseCache <- true
			logging.Err().Alert().Message("HTTP GET request error retryCount exceeded for: /online. Is server up? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.").Add()
		}
		hadError = true
	}

	err = json.NewDecoder(req.Body).Decode(&Cache.Online.Online)
	if err != nil {
		if retryCount == 5 {
			logging.Err().Alert().Message("JSON Unmarshal error retryCount exceeded for: /online. Are API definitions up to date? Cache refresh paused until heartbeat detected. Further charts will be generated with last server snapshot.")
		}
		hadError = true
	}

	Cache.Online.Timestamp = time.Now()

	defer func() {
		if hadError {
			mult := retryCount + 1
			baseTime := time.Duration(30 * mult)
			tick := time.NewTimer(baseTime * time.Second)
			<-tick.C
			syncKills(retryCount + 1)
			tick.Stop()
		}
	}()
}
