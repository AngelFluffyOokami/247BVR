package globaldb

import "gorm.io/gorm"

// Massive fuck you to readability. But necessary in case a function tries to use the DB at same time as another... Damned be SQLite.
var GetDB = make(chan *gorm.DB)
var DoneDB = make(chan bool)

// Workaround for SQLite3. It only lets a single function at a time use the database.
func DBLoop(DB *gorm.DB) {
	for {
		GetDB <- DB
		<-DoneDB
	}
}
