//go:build windows

package win32

import (
	"fmt"

	wapi "github.com/iamacarpet/go-win64api"
	"github.com/tjarratt/babble"
)

func AddUser() bool {

	userOptions := wapi.UserAddOptions{
		Username:  "247bvr",
		Password:  newPassword(),
		FullName:  "HSVR ELOBot Statistics Service User",
		PrivLevel: wapi.USER_PRIV_ADMIN,
		Comment:   "User installed and used to run HSVR ELO statistics bot service.",
	}

	ok, err := wapi.UserAddEx(userOptions)
	if err != nil {

		fmt.Println(err.Error())
		return false
	}
	if !ok {
		return false
	}

	return true
}

func AddService() bool {

	return true
}
func newPassword() string {
	babbler := babble.NewBabbler()
	babbler.Count = 4
	babbler.Separator = "-"
	return babbler.Babble()
}
