package cubecounterimage

import (
	"cci_grapher/db"
	"cci_grapher/utils"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (ccR *cubeCounterRequest) createData(db *db.DataBase) *cubeCounterData {
	usernames, e := db.GetAllUsernames()
	if e != nil {
		utils.ERROR("An error occurred trying to prepare the username database."+e.Error(), "CubeCounter.createData")
		return nil
	}
	now := time.Now()

	var out cubeCounterData = GetCubeCounterDate()
	for _, c := range ccR.channelIDs {
		e := ccR.processDB(c, usernames, &out, db)
		if e != nil {
			utils.ERROR("An error occurred trying to process the database for channel "+c, "CubeCounter.createData")
			return nil
		}
	}
	utils.LOGGING(fmt.Sprintf("Making cubeCounterData took: %v", time.Since(now)), "CCI.createData")
	return &out
}

func (ccr *cubeCounterRequest) processDB(channelID string, userGetter map[string]string, ccB *cubeCounterData, db *db.DataBase) error {
	var activeMembers = map[string]ActiveMembersStruct{}
	rowsStart := time.Now()
	var lastID string = "0"
	var rowsCounter int = 0
	var chunkCounter int = 0
	var countedRows int = 0

	for lastID != "" {
		chunkCounter++
		rows, e := db.GetAllMessagesBetweenForChannelFromID(ccr.startDate, ccr.endDate, channelID, lastID)
		if e != nil {
			utils.ERROR("An error occurred trying to fetch data from "+channelID+"\n"+e.Error(), "CubeCounter.createImg")
			return e
		}
		lastID, countedRows, e = ccB.processRows(rows, userGetter, activeMembers)
		if e != nil {
			utils.ERROR("An error occurred trying to process the rows."+e.Error(), "CubeCounter.processDB")
			return e
		}
		rowsCounter += countedRows
	}

	rowsEnd := time.Now()
	utils.LOGGING(fmt.Sprintf("[%s] Looping over %d rows in %d chunks took: %v", channelID, rowsCounter, chunkCounter, rowsEnd.Sub(rowsStart)), "CCI.processDB")
	return nil
}
