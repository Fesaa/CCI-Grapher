package db

import (
	"cci_grapher/logging"
	"cci_grapher/config"
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB
var GetUsernames *sql.Stmt
var GetAllMessages *sql.Stmt
var GetMessagesForChannel *sql.Stmt
var GetAllMessagesBetween *sql.Stmt
var GetAllMessagesBetweenForChannel *sql.Stmt

func Connect() bool {
	logging.INFO("Connecting to database", "db.Connect")
	var err error
	db, err = sql.Open("postgres", config.CCConfig.PsqlLink)
	err2 := db.Ping()
	if err != nil && err2 != nil {
		if err != nil {
			logging.ERROR(err.Error(), "server")
		} else if err2 != nil {
			logging.ERROR(err2.Error(), "server")
		}
		return false
	} else {
		logging.SUCCESS("Connected to database", "server")
		return true
	}
}

func Disconnect() {
	logging.INFO("Disconnecting from database", "server")
	db.Close()
	logging.SUCCESS("Disconnected from database", "server")
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
	Connect()
	var e error
	GetUsernames, e = Prepare("SELECT user_id,username FROM usernames;")
	if e != nil {
		logging.ERROR(e.Error(), "db.Init")
	}
	GetAllMessages, e = Prepare("SELECT * FROM messages;")
	if e != nil {
		logging.ERROR(e.Error(), "db.Init")
	}
	GetMessagesForChannel, e = Prepare("SELECT * FROM messages WHERE channel_id = $1;")
	if e != nil {
		logging.ERROR(e.Error(), "db.Init")
	}
	GetAllMessagesBetween, e = Prepare("SELECT * FROM messages WHERE time BETWEEN $1 AND $2;")
	if e != nil {
		logging.ERROR(e.Error(), "db.Init")
	}
	GetAllMessagesBetweenForChannel, e = Prepare("SELECT * FROM messages WHERE time BETWEEN $1 AND $2 AND channel_id = $3;")
	if e != nil {
		logging.ERROR(e.Error(), "db.Init")
	}
}