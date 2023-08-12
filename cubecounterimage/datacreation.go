package cubecounterimage

import (
	"cci_grapher/db"
	"cci_grapher/utils"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func processDB(channelID string, ccB *cubeCounterData, ccr cubeCounterRequest, userGetter map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	now := time.Now()
	rows, e := db.GetAllMessagesBetweenForChannel(ccr.startDate, ccr.endDate, channelID)
	if e != nil {
		utils.ERROR("An error occurred trying to fetch data from "+channelID+"\n"+e.Error(), "CubeCounter.createImg")
		return
	}
	utils.LOGGING(fmt.Sprintf("Getting data from channel_db took: %v", time.Since(now)), "CCI.processDB")

	var activeMembers = map[string]ActiveMembersStruct{}

	var totalMessagesTimer time.Duration = 0
	var consecutiveTimeTimer1 time.Duration = 0
	var consecutiveTimeTimer2 time.Duration = 0
	var consecutiveTimeTimer3 time.Duration = 0
	var roleDistributionTimer1 time.Duration = 0
	var roleDistributionTimer2 time.Duration = 0
	for rows.Next() {
		var messageID string
		var channelID string
		var userID string
		var tString string
		var rolesString string
		err := rows.Scan(&messageID, &channelID, &userID, &rolesString, &tString)
		if err != nil {
			utils.ERROR("An error occurred trying to scan from rows."+err.Error(), "CubeCounter.processDB")
			return
		}

		t, err := time.Parse("2006-01-02T15:04:05.999Z", strings.TrimSuffix(strings.Split(tString, "+")[0], " "))
		if err != nil {
			utils.ERROR("Error parsing time;\n "+err.Error(), "CubeCounter.processDB")
			continue
		}

		var userName = userGetter[userID]

		msg := MessageEntry{
			Date:     t,
			AuthorID: userName,
			RolesIDs: strings.Split(rolesString, ","),
		}

		ccB.totalMessageCount = ccB.totalMessageCount + 1

		// Adding to totalMessages
		now = time.Now()
		if _, ok := ccB.totalMessages[msg.AuthorID]; ok {
			ccB.totalMessages[msg.AuthorID] = ccB.totalMessages[msg.AuthorID] + 1
		} else {
			ccB.totalMessages[msg.AuthorID] = 1
		}
		totalMessagesTimer += time.Since(now)

		// consecutiveTime calculations
		var toRemove []string

		now = time.Now()
		if _, ok := activeMembers[msg.AuthorID]; ok {
			temp := activeMembers[msg.AuthorID]
			temp.lastTime = msg.Date
			temp.messages++
			activeMembers[msg.AuthorID] = temp
		} else {
			activeMembers[msg.AuthorID] = ActiveMembersStruct{
				startTime: msg.Date,
				lastTime:  msg.Date,
				messages:  1,
			}
		}
		consecutiveTimeTimer1 += time.Since(now)

		now = time.Now()
		for userID, info := range activeMembers {
			timeDifference := msg.Date.Sub(info.lastTime)

			if timeDifference.Seconds()/60 > 10 {
				timeDelta := info.lastTime.Sub(info.startTime)

				if timeDelta.Seconds() == 0 {
					timeDelta = time.Minute * 2
				}
				if info.messages*10 < int(timeDelta.Seconds()/60) {
					timeDelta = time.Duration(int64(time.Minute) * 2 * int64(info.messages))
				}

				ccB.consecutiveTime[userID] = append(ccB.consecutiveTime[userID], timeDelta.Seconds())

				toRemove = append(toRemove, userID)
			}
		}
		consecutiveTimeTimer2 += time.Since(now)

		now = time.Now()
		for _, userID := range toRemove {
			delete(activeMembers, userID)
		}
		consecutiveTimeTimer3 += time.Since(now)

		// roleDistribution calculations
		var toAdd []string
		var containsJava bool = false
		var containsBedrock bool = false

		now = time.Now()
		for _, role := range msg.RolesIDs {
			if _, ok := roles[role]; ok {
				toAdd = append(toAdd, role)
				ccB.roleDistribution[role] = ccB.roleDistribution[role] + 1
			}
			if role == javaRole {
				containsJava = true
			}
			if role == bedrockRole {
				containsBedrock = true
			}
		}
		roleDistributionTimer1 += time.Since(now)

		now = time.Now()
		if containsJava && !containsBedrock {
			ccB.roleDistribution["Java Only"] = ccB.roleDistribution["Java Only"] + 1
		} else if !containsJava && containsBedrock {
			ccB.roleDistribution["Bedrock Only"] = ccB.roleDistribution["Bedrock Only"] + 1
		} else if containsJava && containsBedrock {
			ccB.roleDistribution["Dual"] = ccB.roleDistribution["Dual"] + 1
		}

		if len(toAdd) == 0 {
			ccB.roleDistribution["No Roles"] = ccB.roleDistribution["No Roles"] + 1
		}
		roleDistributionTimer2 += time.Since(now)

		// hourlyActivity
		ccB.hourlyActivity[msg.Date.Hour()] = ccB.hourlyActivity[msg.Date.Hour()] + 1
	}

	utils.LOGGING(fmt.Sprintf("[%s] Counting total messages per person took: %v", channelID, totalMessagesTimer), "CCI.processDB")
	utils.LOGGING(fmt.Sprintf("[%s] Adding or updating consecutive struct took: %v", channelID, consecutiveTimeTimer1), "CCI.processDB")
	utils.LOGGING(fmt.Sprintf("[%s] Looping over active members took: %v", channelID, consecutiveTimeTimer2), "CCI.processDB")
	utils.LOGGING(fmt.Sprintf("[%s] Removing inactive members took: %v", channelID, consecutiveTimeTimer3), "CCI.processDB")
	utils.LOGGING(fmt.Sprintf("[%s] Looping over user roles took: %v", channelID, roleDistributionTimer1), "CCI.processDB")
	utils.LOGGING(fmt.Sprintf("[%s] Checking special roles took: %v", channelID, roleDistributionTimer1), "CCI.processDB")

}

func createData(ccR cubeCounterRequest) *cubeCounterData {
	now := time.Now()
	data, e := db.GetAllUsernames()
	if e != nil {
		utils.ERROR("An error occurred trying to prepare the username database."+e.Error(), "CubeCounter.createData")
		return &cubeCounterData{}
	}

	usernames := make(map[string]string)
	for data.Next() {
		var username string
		var userId string
		err := data.Scan(&userId, &username)
		if err != nil {
			utils.ERROR("An error occurred trying to scan from data", "CubeCounter.createData")
			return nil
		}
		usernames[userId] = username
	}
	utils.LOGGING(fmt.Sprintf("Making usernames map took: %v", time.Since(now)), "CCI.createData")

	now = time.Now()
	var ccB = cubeCounterData{
		totalMessageCount: 0,
		totalMessages:     make(map[string]int),
		consecutiveTime:   make(map[string][]float64),
		roleDistribution:  make(map[string]int),
		hourlyActivity:    make(map[int]int),
	}

	for k := range roles {
		ccB.roleDistribution[k] = 0
	}

	for hour := 0; hour < 24; hour++ {
		ccB.hourlyActivity[hour] = 0
	}
	utils.LOGGING(fmt.Sprintf("Making cubeCounterData took: %v", time.Since(now)), "CCI.createData")

	now = time.Now()
	wg := sync.WaitGroup{}
	wg.Add(len(ccR.channelIDs))
	for _, c := range ccR.channelIDs {
		utils.INFO(fmt.Sprintf("Starting `processDB` for channelID: \"%s\"", c), "CubeCounter.createData")
		go processDB(c, &ccB, ccR, usernames, &wg)
	}
	utils.LOGGING(fmt.Sprintf("processDB took: %v", time.Since(now)), "CCI.createData")

	wg.Wait()
	return &ccB
}
