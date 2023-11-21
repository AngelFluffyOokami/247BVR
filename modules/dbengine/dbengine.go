package dbengine

import (
	"fmt"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	scribble "github.com/nanobox-io/golang-scribble"
)

// DBv todo
var DBv = &DB{}

// DB Struct holds values and functions for the engine.
type DB struct {
	Db *scribble.Driver
}

// Init Initializes database engine.
func (ctx *DB) Init() {
	ctx.initDB()

}

func (ctx *DB) initDB() {

	var err error
	// a new scribble driver, providing the directory where it will be writing to,
	// and a qualified logger if desired
	ctx.Db, err = scribble.New("./vtol.vr", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}

// WriteDB todo
func (ctx *DB) WriteDB(dataType string, data any, pid string) {

	var quickAsserted bool

	switch dataType {
	case "kill":
		killsData, ok := data.(global.KillEvent)
		if !ok {
			fmt.Println("bad type assertion, breaking from QuickAssert and attempting to assert proper data type.")
			quickAsserted = false
			break
		}
		quickAsserted = true
		ctx.Db.Write("kill", pid, killsData)

	case "online":
		onlineData, ok := data.([]global.WsOnlineData)

		if !ok {
			fmt.Println("bad type assertion, breaking from QuickAssert and attempting to assert proper data type.")
			quickAsserted = false
			break
		}
		quickAsserted = true
		ctx.Db.Write("online", pid, onlineData)
	case "spawn":
		spawnData, ok := data.(global.WsSpawnData)

		if !ok {
			fmt.Println("bad type assertion, breaking from QuickAssert and attempting to assert proper data type.")
			quickAsserted = false
			break
		}
		quickAsserted = true
		ctx.Db.Write("spawn", pid, spawnData)
	case "login":
		loginData, ok := data.(global.WsUserLogEvent)
		if !ok {
			fmt.Println("bad type assertion, breaking from QuickAssert and attempting to assert proper data type.")
			quickAsserted = false
			break
		}
		quickAsserted = true
		ctx.Db.Write("login", pid, loginData)
	case "logout":
		logoutData, ok := data.(global.WsUserLogEvent)
		if !ok {
			fmt.Println("bad type assertion, breaking from QuickAssert and attempting to assert proper data type.")
			quickAsserted = false
			break
		}
		quickAsserted = true
		ctx.Db.Write("login", pid, logoutData)
	case "users":
		userData, ok := data.(global.User)
		if !ok {
			// this log message is ostensably a lie judging from how it never really tries doing anytthing if QuickAssert fails.
			fmt.Println("bad type assertion, breaking from QuickAssert and attempting to assert proper data type.")
			quickAsserted = false
			break
		}
		quickAsserted = true
		ctx.Db.Write("users", pid, userData)
	default:
		ctx.Db.Write(dataType, pid, data.([]string))
	}

	if !quickAsserted {
		fmt.Println("false")
	}
}
