package test

import (
	"fmt"

	"github.com/angelfluffyookami/247BVR/modules/bvr"
	"github.com/angelfluffyookami/247BVR/modules/common/global"
	"github.com/google/go-cmp/cmp"
)

func Test() {

	var User1 bvr.UserStruct
	var User2 bvr.UserStruct
	var User3 bvr.UserStruct
	var User4 bvr.UserStruct
	for _, x := range bvr.Cache.Users.Users {

		switch x.ID {
		case global.TestID1:
			User1 = x
		case global.TestID2:
			User2 = x
		case global.TestID3:
			User3 = x
		case global.TestID4:
			User4 = x
		}

		if !cmp.Equal(User1, bvr.UserStruct{}) && !cmp.Equal(User2, bvr.UserStruct{}) && !cmp.Equal(User3, bvr.UserStruct{}) && cmp.Equal(User4, bvr.UserStruct{}) {
			break
		}
	}

	fmt.Println("U1LI: " + fmt.Sprint(len(User1.LoginTimes)) + " U1LO: " + fmt.Sprint(len(User1.LogoutTimes)) + `
	` + "U2LI: " + fmt.Sprint(len(User2.LoginTimes)) + " U2LO: " + fmt.Sprint(len(User2.LogoutTimes)) + `
	` + "U3LI: " + fmt.Sprint(len(User3.LoginTimes)) + " U3LO: " + fmt.Sprint(len(User3.LogoutTimes)) + `
	` + "U4LI: " + fmt.Sprint(len(User4.LoginTimes)) + " U4LO: " + fmt.Sprint(len(User4.LogoutTimes)))

}
