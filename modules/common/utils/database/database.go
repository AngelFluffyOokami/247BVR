package database

import (
	"github.com/angelfluffyookami/247BVR/modules/common/utils/database/coredb"
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
