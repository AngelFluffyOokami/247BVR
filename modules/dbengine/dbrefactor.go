package dbengine

import (
	"fmt"
	"os"
	"runtime"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	scribble "github.com/nanobox-io/golang-scribble"
)

var basePath string

type User struct {
	UserID string
	User   global.User
	lookup bool
	db     *scribble.Driver
	dir    string
}

func (ctx *User) AddOrUpdateUser() {

	if !ctx.lookup {
		ctx.lookupUser()
	}
}

func (ctx *User) lookupUser() {

	if _, err := os.Stat(basePath + "/247BVR/" + ctx.UserID); os.IsNotExist(err) {
		// path/to/whatever does not exist
	}
	ctx.lookup = true
}

func setBasePath() {
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("C:\\Users\\Public\\247BVR"); err != nil {
			fmt.Println("247BVR folder does not exists, attempting to create.")
			os.Mkdir("C:\\Users\\Public\\247BVR", 0755)
		}

	}
}
