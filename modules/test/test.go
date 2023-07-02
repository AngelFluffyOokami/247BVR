package test

import (
	"fmt"
	"sort"

	"github.com/angelfluffyookami/247BVR/modules/bvr"
	"github.com/angelfluffyookami/247BVR/modules/common/global"
)

func Test() {

	User1, err := bvr.GetUserByID(global.TestID1)

	if err != nil {
		return
	}

	User2, err := bvr.GetUserByID(global.TestID2)

	if err != nil {
		return
	}

	User3, err := bvr.GetUserByID(global.TestID3)

	if err != nil {
		return
	}

	User4, err := bvr.GetUserByID(global.TestID4)

	if err != nil {
		return
	}

	var User1Kills []bvr.KillStruct
	for _, x := range bvr.Cache.Kills.Kills {

		if x.KillerID == User1.ID {

			User1Kills = append(User1Kills, x)

		}

	}

	var User1Deaths []bvr.DeathStruct

	for _, x := range bvr.Cache.Deaths.Deaths {

		if x.VictimID == User1.ID {
			User1Deaths = append(User1Deaths, x)
		}
	}

	var UserActivity []int64
	for _, x := range User1Kills {
		UserActivity = append(UserActivity, x.Time)
	}
	for _, x := range User1Deaths {
		UserActivity = append(UserActivity, x.Time)
	}

	fmt.Println("U1LI: " + fmt.Sprint(len(User1.LoginTimes)) + " U1LO: " + fmt.Sprint(len(User1.LogoutTimes)) + `
` + "U2LI: " + fmt.Sprint(len(User2.LoginTimes)) + " U2LO: " + fmt.Sprint(len(User2.LogoutTimes)) + `
` + "U3LI: " + fmt.Sprint(len(User3.LoginTimes)) + " U3LO: " + fmt.Sprint(len(User3.LogoutTimes)) + `
` + "U4LI: " + fmt.Sprint(len(User4.LoginTimes)) + " U4LO: " + fmt.Sprint(len(User4.LogoutTimes)))

	var logs logslice

	var logins []int64
	var logouts []int64

	logs = append(logs, User1.LoginTimes...)
	logs = append(logs, User1.LogoutTimes...)
	logins = append(logins, User1.LoginTimes...)
	logouts = append(logouts, User1.LogoutTimes...)

	sort.Sort(logs)

	var logType []loginout

	for _, x := range logs {
		found := false
		for _, y := range User1.LoginTimes {
			if x == y {

				logType = append(logType, loginout{timestamp: x, login: true})
				found = true
				break
			}
		}

		if found {
			continue
		}
		for _, y := range User1.LogoutTimes {
			if x == y {
				logType = append(logType, loginout{timestamp: x, login: false})
				break
			}
		}
	}

	var missingPairIndex []int

	for x, y := range logType {
		if x == 0 {
			continue
		}

		if y.login == logType[x-1].login {

			switch y.login {
			case true:
				fmt.Println(fmt.Sprint(x-1) + " missing logout")
				missingPairIndex = append(missingPairIndex, x-1)

				for q, z := range UserActivity {
					if q == 0 {
						continue
					}
					if z >= logType[x-1].timestamp {
						logouts = append(logouts, UserActivity[q-1])
						break
					}
				}

			case false:
				fmt.Println(fmt.Sprint(x-1) + " missing login")
				missingPairIndex = append(missingPairIndex, x-1)
				for _, z := range UserActivity {

					if z >= logType[x-1].timestamp {
						logins = append(logins, z)
						break
					}
				}
			}

		} else {
			continue
		}
	}

	logs = logins
	logs = append(logs, logouts...)

	sort.Sort(logs)

	logType = []loginout{}

	for _, x := range logs {
		found := false
		for _, y := range User1.LoginTimes {
			if x == y {

				logType = append(logType, loginout{timestamp: x, login: true})
				found = true
				break
			}
		}

		if found {
			continue
		}
		for _, y := range User1.LogoutTimes {
			if x == y {
				logType = append(logType, loginout{timestamp: x, login: false})
				break
			}
		}
	}

	for x, y := range logType {
		if x == 0 {
			continue
		}
		if y.login == logType[x-1].login {

			switch y.login {
			case true:
				fmt.Println(fmt.Sprint(x-1) + " missing logout")
				missingPairIndex = append(missingPairIndex, x-1)

			case false:
				fmt.Println(fmt.Sprint(x-1) + " missing login")
				missingPairIndex = append(missingPairIndex, x-1)
			}
		} else {
			continue
		}
	}

	fmt.Println(fmt.Sprint(len(logType)))

	fmt.Println("noh")

}

type loginout struct {
	timestamp int64
	login     bool
}

type logslice []int64

func (a logslice) Len() int           { return len(a) }
func (a logslice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a logslice) Less(i, j int) bool { return a[i] < a[j] }
