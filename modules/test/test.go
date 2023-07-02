package test

import (
	"fmt"

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

	fmt.Println("U1LI: " + fmt.Sprint(len(User1.LoginTimes)) + " U1LO: " + fmt.Sprint(len(User1.LogoutTimes)) + `
	` + "U2LI: " + fmt.Sprint(len(User2.LoginTimes)) + " U2LO: " + fmt.Sprint(len(User2.LogoutTimes)) + `
	` + "U3LI: " + fmt.Sprint(len(User3.LoginTimes)) + " U3LO: " + fmt.Sprint(len(User3.LogoutTimes)) + `
	` + "U4LI: " + fmt.Sprint(len(User4.LoginTimes)) + " U4LO: " + fmt.Sprint(len(User4.LogoutTimes)))

}
