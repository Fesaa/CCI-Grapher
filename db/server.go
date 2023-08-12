package db

import (
	"cci_grapher/config"
	"cci_grapher/utils"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var getUsernames *sql.Stmt
var getAllMessagesBetweenForChannel *sql.Stmt

func Connect() bool {
	utils.INFO("Connecting to database", "db.Connect")
	var err error
	db, err = sql.Open("postgres", config.CCConfig.PsqlLink)
	err2 := db.Ping()
	if err != nil || err2 != nil {
		if err != nil {
			utils.ERROR(err.Error(), "server")
		} else if err2 != nil {
			utils.ERROR(err2.Error(), "server")
		}
		return false
	} else {
		utils.SUCCESS("Connected to database", "server")
		return true
	}
}

func Disconnect() error {
	utils.INFO("Disconnecting from database", "server")
	err := db.Close()
	if err != nil {
		return err
	}
	utils.SUCCESS("Disconnected from database", "server")
	return nil
}

func Prepare(q string) (*sql.Stmt, error) {
	r, e := db.Prepare(q)
	if e == nil {
		return r, e
	} else {
		return nil, e
	}
}

func Init() {
	var e error
	getUsernames, e = Prepare("SELECT user_id,username FROM usernames;")
	if e != nil {
		utils.ERROR(e.Error(), "db.Init")
	}
	getAllMessagesBetweenForChannel, e = Prepare("SELECT * FROM messages WHERE time BETWEEN $1 AND $2 AND channel_id = $3;")
	if e != nil {
		utils.ERROR(e.Error(), "db.Init")
	}
}

func GetAllMessagesBetweenForChannel(start time.Time, end time.Time, channelId string) (*sql.Rows, error) {
	rows, err := getAllMessagesBetweenForChannel.Query(start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), channelId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetAllUsernames() (*sql.Rows, error) {
	return getUsernames.Query()
}
