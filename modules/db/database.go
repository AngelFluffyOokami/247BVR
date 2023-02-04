package database

import (
	coredb "github.com/angelfluffyookami/247BVR/modules/db/core_models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	cinnamondb, err := gorm.Open(sqlite.Open("database/cinnamon.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	cinnamondb.AutoMigrate(&coredb.Cinnamon{}, coredb.Kill{}, coredb.User{}, coredb.LimitedUserData{}, coredb.Death{})

	return cinnamondb
}
