package dbengine

import (
	"fmt"
	"os"
	"runtime"

	"github.com/angelfluffyookami/HSVRUSB/modules/common/global"
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

	if _, err := os.Stat(basePath + "/HSVRUSB/" + ctx.UserID); os.IsNotExist(err) {
		// path/to/whatever does not exist
	}
	ctx.lookup = true
}

func setBasePath() {
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("C:\\Users\\Public\\HSVRUSB"); err != nil {
			fmt.Println("HSVRUSB folder does not exists, attempting to create.")
			os.Mkdir("C:\\Users\\Public\\HSVRUSB", 0755)
		}

	}
}
