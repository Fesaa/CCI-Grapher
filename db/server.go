package db

import (
	"cci_grapher/utils"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type DataBase struct {
	db *sql.DB
	getUsernames *sql.Stmt
	getAllMessagesBetweenForChannel *sql.Stmt
    getAllMessagesBetweenForChannelForUsers *sql.Stmt
	getAllMessagesBetweenForChannelFromID *sql.Stmt
    getAllMessagesBetweenForChannelFromIDForUsers *sql.Stmt
}

func (d *DataBase) Connect(psql string) error {
	utils.INFO("Connecting to database", "server.Connect")

	var err error
	d.db, err = sql.Open("postgres", psql)
	if err != nil {
		return err
	}
	err = d.db.Ping()
	if err != nil {
		return err
	}
	utils.SUCCESS("Connected to database", "server.Connect")
	return nil
}

func (d *DataBase) Disconnect() error {
	utils.INFO("Disconnecting from database", "server.Disconnect")
	err := d.db.Close()
	if err != nil {
		return err
	}
	utils.SUCCESS("Disconnected from database", "server.Disconnect")
	return nil
}

func (d *DataBase) Init() error {
	var e error
	d.getUsernames, e = d.db.Prepare("SELECT user_id,username FROM usernames;")
	if e != nil {
		return e
	}
	d.getAllMessagesBetweenForChannel, e = d.db.Prepare("SELECT message_id,user_id,roles,time FROM messages WHERE time BETWEEN $1 AND $2 AND channel_id = $3 LIMIT 10000;")
	if e != nil {
		return e
	}
    d.getAllMessagesBetweenForChannelForUsers, e = d.db.Prepare("SELECT message_id,user_id,roles,time FROM messages WHERE time BETWEEN $1 AND $2 AND channel_id = $3 AND user_id = ANY($4) LIMIT 10000;")
    if e != nil {
        return e
    }
	d.getAllMessagesBetweenForChannelFromID, e = d.db.Prepare("SELECT message_id,user_id,roles,time FROM messages WHERE time BETWEEN $1 AND $2 AND channel_id = $3 AND message_id > $4 LIMIT 10000;")
	if e != nil {
		return e
	}
    d.getAllMessagesBetweenForChannelFromIDForUsers, e = d.db.Prepare("SELECT message_id,user_id,roles,time FROM messages WHERE time BETWEEN $1 AND $2 AND channel_id = $3 AND message_id > $4 AND user_id = ANY($5) LIMIT 10000;")
    if e != nil {
        return e
    }

	return nil
}

func (d *DataBase) GetAllMessagesBetweenForChannel(start time.Time, end time.Time, channelId string, userIDs []string) (*sql.Rows, error) {
    var rows *sql.Rows
    var err error
    if len(userIDs) != 0 {
        rows, err = d.getAllMessagesBetweenForChannelForUsers.Query(start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), channelId, pq.Array(userIDs))
    } else {
	    rows, err = d.getAllMessagesBetweenForChannel.Query(start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), channelId)
    }
    if err != nil {
        return nil, err
    }
	return rows, nil
}

func (d *DataBase) GetAllMessagesBetweenForChannelFromID(start time.Time, end time.Time, channelId string, messageId string, userIDs []string) (*sql.Rows, error) {
    var rows *sql.Rows
    var err error
    if len(userIDs) != 0 {
        rows, err = d.getAllMessagesBetweenForChannelFromIDForUsers.Query(start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), channelId, messageId, pq.Array(userIDs))
    } else {
	    rows, err = d.getAllMessagesBetweenForChannelFromID.Query(start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), channelId, messageId)
    }
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DataBase) GetAllUsernames() (map[string]string, error) {
	now := time.Now()
	data, e := d.getUsernames.Query()
	if e != nil {
		utils.ERROR("An error occurred trying to prepare the username database."+e.Error(), "server.GetAllUsernames")
		return nil, e
	}

	usernames := make(map[string]string)
	for data.Next() {
		var username string
		var userId string
		err := data.Scan(&userId, &username)
		if err != nil {
			utils.ERROR("An error occurred trying to scan from data", "server.GetAllUsernames")
			return nil, e
		}
		usernames[userId] = username
	}
	utils.LOGGING(fmt.Sprintf("Making usernames map took: %v", time.Since(now)), "server.GetAllUsernames")
	return usernames, nil
}
