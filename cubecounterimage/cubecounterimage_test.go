package cubecounterimage

import (
	"cci_grapher/config"
	"cci_grapher/db"
	"cci_grapher/utils"
	"testing"
	"time"
)

func TestDataCreation(t *testing.T) {
	config.LoadConfig("../config.json")

	db := db.DataBase{}
	if db.Connect(config.Config.PsqlLink) != nil {
		t.Fail()
		return
	}
	db.Init()

	now := time.Now()
	var ccr cubeCounterRequest = cubeCounterRequest{
		channelIDs: []string{"174837853778345984", "493739733147451402", "725655075069886474"},
		startDate:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(time.Hour * 24 * -7),
		endDate:    now,
	}
	usernames, e := db.GetAllUsernames()
	if e != nil {
		utils.ERROR("An error occurred trying to prepare the username database."+e.Error(), "CubeCounter.createData")
		return
	}

	ccB := GetCubeCounterDate()
	for _, c := range ccr.channelIDs {
		ccr.processDB(c, usernames, &ccB, &db)
	}

	utils.SUCCESS("Finished processing database in: "+time.Since(now).String(), "CubeCounter.processDB")
}
