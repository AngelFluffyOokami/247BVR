package dbengine

import (
	"encoding/json"
	"fmt"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	scribble "github.com/nanobox-io/golang-scribble"
)

// DBv todo
var DBv = &DB{}

// DB Struct holds values and functions for the engine.
type DB struct {
	Kills *scribble.Driver
}

// Init Initializes database engine.
func (ctx DB) Init() {
	ctx.initKillDB()
}

func (ctx DB) initKillDB() {

	var err error
	// a new scribble driver, providing the directory where it will be writing to,
	// and a qualified logger if desired
	ctx.Kills, err = scribble.New("./kills.vtol.vr", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}

// WriteDB todo
func (ctx DB) WriteDB(dataType int, data any) {

	switch dataType {
	case 0:
		var KillsData global.KillData
		err := json.Unmarshal(data.([]byte), &KillsData)
		if err != nil {
			fmt.Println("bad unmarshal: ", err, "\n", "skipping QuickAssert, and trying to find proper data type.")

		}
	}
}
