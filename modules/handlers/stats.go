package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/angelfluffyookami/HSVRUSB/modules/common/global"
	"github.com/angelfluffyookami/HSVRUSB/modules/dbengine"
)

var kills global.Kills

var users global.Users

var deaths global.Deaths

func populateUserVar() error {
	userstring, err := dbengine.DBv.Db.ReadAll("users")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, v := range userstring {
		var userjson global.User
		err := json.Unmarshal([]byte(v), &userjson)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		users = append(users, userjson)
	}

	if users != nil {
		return nil
	} else {
		return fmt.Errorf("uservar nil, unsure why")
	}
}

func populateKillVar() error {
	killstring, err := dbengine.DBv.Db.ReadAll("kill")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, v := range killstring {
		var killjson global.KillEvent
		err := json.Unmarshal([]byte(v), &killjson)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		kills = append(kills, killjson)
	}

	if users != nil {
		return nil
	} else {
		return fmt.Errorf("uservar nil, unsure why")
	}
}

func checkPopulate() bool {
	if users == nil {
		err := populateUserVar()
		if err != nil {
			fmt.Println("error, cannot continue, returning.")
			return false
		}
	}

	if kills == nil {

		err := populateKillVar()
		if err != nil {
			fmt.Println("error, cannot continue, returning.")
			return false
		}

	}

	return true
}

func LowestElo() {
	lowestvalue := 99999999.1

	var lowestUsers global.Users
	for _, v := range users {
		for _, x := range v.EloHistory {
			if x.Elo < lowestvalue {
				lowestvalue = x.Elo
				lowestUsers = nil
				lowestUsers = append(lowestUsers, v)
			} else if x.Elo == lowestvalue {
				donotappend := false
				for _, z := range lowestUsers {
					if z.ID0 == v.ID0 {
						donotappend = true
					}
				}
				if !donotappend {
					lowestUsers = append(lowestUsers, v)
				}

			}
		}
	}

	var usernames string
	for _, v := range lowestUsers {
		for _, x := range v.PilotNames {
			usernames = usernames + "/" + x
		}
		usernames = usernames + "|" + v.ID0 + ", "
	}
	fmt.Println("Lowest elo ever reached: " + fmt.Sprint(lowestvalue) + "\nBy:" + usernames)

}
