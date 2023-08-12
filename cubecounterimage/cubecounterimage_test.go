package cubecounterimage

import (
	"cci_grapher/config"
	"cci_grapher/db"
	"cci_grapher/utils"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDataCreation(t *testing.T) {
	config.LoadConfig("../config.json")
	if !db.Connect() {
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
	data, e := db.GetAllUsernames()
	if e != nil {
		utils.ERROR("An error occurred trying to prepare the username database."+e.Error(), "CubeCounter.createData")
		return
	}

	usernames := make(map[string]string)
	for data.Next() {
		var username string
		var userId string
		err := data.Scan(&userId, &username)
		if err != nil {
			utils.ERROR("An error occurred trying to scan from data", "CubeCounter.createData")
			return
		}
		usernames[userId] = username
	}
	utils.LOGGING(fmt.Sprintf("Making usernames map took: %v", time.Since(now)), "CCI.createData")

	wg := sync.WaitGroup{}
	wg.Add(len(ccr.channelIDs))
	ch := make(chan cubeCounterData, len(ccr.channelIDs))
	for _, c := range ccr.channelIDs {
		go processDB(c, ccr, usernames, &wg, ch)
	}
	wg.Wait()
	close(ch)
	utils.SUCCESS("Finished processing database in: "+time.Since(now).String(), "CubeCounter.processDB")
}
