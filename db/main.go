package db

import (
	"cci_grapher/config"
	"cci_grapher/logging"
	"database/sql"
)

var USERNAMEDB *sql.DB
var CHANNELDBS map[string]*sql.DB = make(map[string]*sql.DB)

func Init() bool {
	for _, channelID := range config.CC.ChannelIDs {
		db, e := sql.Open("sqlite3", config.CC.CubePath+channelID+".sql")
		if e != nil {
			logging.ERROR("An error occurred trying to open the database for channel with id: "+channelID+"\n"+e.Error(), "CubeCounter.createImg")
			return false
		}
		CHANNELDBS[channelID] = db
	}

	usernamesDB, e := sql.Open("sqlite3", config.CC.CubePath+"usernames.sql")
	if e != nil {
		logging.ERROR("An error occurred trying to open the username database."+e.Error(), "CubeCounter.createImg")
		return false
	}
	USERNAMEDB = usernamesDB

	return true
}

func Shutdown() {
	err := USERNAMEDB.Close()
	if err != nil {
		logging.ERROR("Could not close USERNAMEDB", "Shutdown")
		return
	}
	for _, db := range CHANNELDBS {
		err := db.Close()
		if err != nil {
			logging.ERROR("Could not close a CHANNELDB", "Shutdown")
			return
		}
	}
}
