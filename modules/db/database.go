package database

import (
	coredb "github.com/angelfluffyookami/247BVR/modules/db/core_models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	cinnamondb, err := gorm.Open(sqlite.Open("butter.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	cinnamondb.AutoMigrate(&coredb.Cinnamon{})

	return cinnamondb
}
