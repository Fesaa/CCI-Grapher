package cubecounterimage

import (
	"cci_grapher/db"
	"cci_grapher/utils"
	"fmt"
	"sync"
	"time"
)

func (ccR *cubeCounterRequest) createData(db *db.DataBase) *cubeCounterData {
	usernames, e := db.GetAllUsernames()
	if e != nil {
		utils.ERROR("An error occurred trying to prepare the username database."+e.Error(), "CubeCounter.createData")
		return nil
	}
	now := time.Now()

	var out cubeCounterData = GetCubeCounterDate()
	wg := sync.WaitGroup{}
	wg.Add(len(ccR.channelIDs))
	ch := make(chan Message)

	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go merger(ch, &out, &wg1)
	for _, c := range ccR.channelIDs {
		go ccR.processDB(c, usernames, db, &wg, ch)
	}

	wg.Wait()
	close(ch)
	wg1.Wait()
	utils.LOGGING(fmt.Sprintf("Making cubeCounterData took: %v", time.Since(now)), "CCI.createData")
	return &out
}

func (ccr *cubeCounterRequest) processDB(channelID string, userGetter map[string]string,
	db *db.DataBase, wg *sync.WaitGroup, ch chan Message) {
	defer wg.Done()
	var activeMembers = map[string]ActiveMembersStruct{}
	rowsStart := time.Now()
	var lastID string = "0"
	var rowsCounter int = 0
	var chunkCounter int = 0
	var countedRows int = 0

	for lastID != "" {
		chunkCounter++
		rows, e := db.GetAllMessagesBetweenForChannelFromID(ccr.startDate, ccr.endDate, channelID, lastID, ccr.userIDs)
		if e != nil {
			utils.ERROR("An error occurred trying to fetch data from "+channelID+"\n"+e.Error(), "CubeCounter.createImg")
			return
		}
		lastID, countedRows, e = processRows(rows, userGetter, activeMembers, ch)
		if e != nil {
			utils.ERROR("An error occurred trying to process the rows."+e.Error(), "CubeCounter.processDB")
			return
		}
		rowsCounter += countedRows
	}

	rowsEnd := time.Now()
	utils.LOGGING(fmt.Sprintf("[%s] Looping over %d rows in %d chunks took: %v", channelID, rowsCounter, chunkCounter, rowsEnd.Sub(rowsStart)), "CCI.processDB")
	if len(ccr.userIDs) == 1 && rowsCounter > 0 {
		ch <- Message{
			cubeCounterDataType: ChannelTotalMessages,
			data:                []interface{}{channelID, rowsCounter},
		}
	}
}

func merger(ch chan Message, out *cubeCounterData, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ch {
		switch msg.cubeCounterDataType {
		case RoleDistribution:
			role := msg.data.(string)
			if _, ok := out.roleDistribution[role]; !ok {
				out.roleDistribution[role] = 1
			} else {
				out.roleDistribution[role]++
			}
		case HourlyActivity:
			hour := msg.data.(int)
			if _, ok := out.hourlyActivity[hour]; !ok {
				out.hourlyActivity[hour] = 1
			} else {
				out.hourlyActivity[hour]++
			}
		case TotalMessageCount:
			out.totalMessageCount++
		case TotalMessages:
			authorID := msg.data.(string)
			if _, ok := out.totalMessages[authorID]; !ok {
				out.totalMessages[authorID] = 1
			} else {
				out.totalMessages[authorID]++
			}
		case ConsecutiveTime:
			authorID := msg.data.([]interface{})[0].(string)
			time := msg.data.([]interface{})[1].(float64)
			if _, ok := out.consecutiveTime[authorID]; !ok {
				out.consecutiveTime[authorID] = []float64{time}
			} else {
				out.consecutiveTime[authorID] = append(out.consecutiveTime[authorID], time)
			}
		case ChannelTotalMessages:
			channelName := channels[msg.data.([]interface{})[0].(string)]
			out.totalMessages[channelName] = msg.data.([]interface{})[1].(int)
		}
	}
}
