package cubecounterimage

import (
	"cci_grapher/db"
	"cci_grapher/utils"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func createData(ccR cubeCounterRequest, db *db.DataBase) *cubeCounterData {
	usernames, e := db.GetAllUsernames()
	if e != nil {
		utils.ERROR("An error occurred trying to prepare the username database."+e.Error(), "CubeCounter.createData")
		return nil
	}
	now := time.Now()

	var out cubeCounterData = GetCubeCounterDate()
	for _, c := range ccR.channelIDs {
		utils.INFO(fmt.Sprintf("Starting `processDB` for channelID: \"%s\"", c), "CubeCounter.createData")
		e := processDB(c, ccR, usernames, &out, db)
		if e != nil {
			utils.ERROR("An error occurred trying to process the database for channel " + c, "CubeCounter.createData")
			return nil
		}
	}
	utils.LOGGING(fmt.Sprintf("Making cubeCounterData took: %v", time.Since(now)), "CCI.createData")
	return &out
}

func processDB(channelID string, ccr cubeCounterRequest, userGetter map[string]string, ccB *cubeCounterData, db *db.DataBase) error {
	now := time.Now()
	rows, e := db.GetAllMessagesBetweenForChannel(ccr.startDate, ccr.endDate, channelID)
	if e != nil {
		utils.ERROR("An error occurred trying to fetch data from "+channelID+"\n"+e.Error(), "CubeCounter.createImg")
		return e
	}
	utils.LOGGING(fmt.Sprintf("Getting data from channel_db took: %v", time.Since(now)), "CCI.processDB")

	var activeMembers = map[string]ActiveMembersStruct{}

	rowsStart := time.Now()
	for rows.Next() {
		var messageID string
		var channelID string
		var userID string
		var tString string
		var rolesString string
		err := rows.Scan(&messageID, &channelID, &userID, &rolesString, &tString)
		if err != nil {
			utils.ERROR("An error occurred trying to scan from rows."+err.Error(), "CubeCounter.processDB")
			return err
		}
		t, err := time.Parse("2006-01-02T15:04:05.999Z", strings.TrimSuffix(strings.Split(tString, "+")[0], " "))
		if err != nil {
			utils.ERROR("Error parsing time;\n "+err.Error(), "CubeCounter.processDB")
			continue
		}

		msg := MessageEntry{
			Date:     t,
			AuthorID: userGetter[userID],
			RolesIDs: strings.Split(rolesString, ","),
		}
		ccB.AddRowInfo(msg, activeMembers)
	}
	rowsEnd := time.Now()
	utils.LOGGING(fmt.Sprintf("[%s] Looping over rows took: %v", channelID, rowsEnd.Sub(rowsStart)), "CCI.processDB")
	utils.INFO(fmt.Sprintf("[%s] Finished processing database, had %d msgs", channelID, ccB.totalMessageCount), "CubeCounter.processDB")
	return nil
}
