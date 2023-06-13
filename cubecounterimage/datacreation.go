package cubecounterimage

import (
	"cci_grapher/config"
	"cci_grapher/logging"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func processDB(channelID string, ccB *cubeCounterData, ccr cubeCounterRequest, userGetter map[string]string) {
	db, e := sql.Open("sqlite3", config.CC.CubePath+channelID+".sql")
	if e != nil {
		logging.ERROR("An error occurred trying to open the database for channel with id: "+channelID+"\n"+e.Error(), "CubeCounter.createImg")
		return
	}
	defer db.Close()

	//ccr.startDate = ccr.startDate.Add(time.Hour * -24)
	//ccr.endDate = ccr.endDate.Add(time.Hour * 24)

	rows, e := db.Query("SELECT user_id,time,roles FROM msgs WHERE time > $1 AND time < $2;", ccr.startDate, ccr.endDate)
	if e != nil {
		logging.ERROR("An error occurred trying to fetch data from "+channelID+"\n"+e.Error(), "CubeCounter.createImg")
		return
	}

	var active_members = map[string]ActiveMembersStruct{}

	var totalMessagesTimer time.Duration = 0
	var consecutiveTimeTimer1 time.Duration = 0
	var consecutiveTimeTimer2 time.Duration = 0
	var consecutiveTimeTimer3 time.Duration = 0
	var roleDistributionTimer1 time.Duration = 0
	var roleDistributionTimer2 time.Duration = 0
	for rows.Next() {
		var userID string
		var tString string
		var rolesString string
		rows.Scan(&userID, &tString, &rolesString)

		t, err := time.Parse("2006-01-02 15:04:05.999", strings.TrimSuffix(strings.Split(tString, "+")[0], " "))
		if err != nil {
			logging.ERROR("Error parsing time;\n "+err.Error(), "CubeCounter.createImg")
			continue
		}

		var userName string = userGetter[userID]

		msg := MessageEntry{
			Date:     t,
			AuthorID: userName,
			RolesIDs: strings.Split(rolesString, ","),
		}

		ccB.totalMessageCount = ccB.totalMessageCount + 1

		// Adding to totalMessages
		now := time.Now()
		if _, ok := ccB.totalMessages[msg.AuthorID]; ok {
			ccB.totalMessages[msg.AuthorID] = ccB.totalMessages[msg.AuthorID] + 1
		} else {
			ccB.totalMessages[msg.AuthorID] = 1
		}
		totalMessagesTimer += time.Since(now)

		// consecutiveTime calculations
		var to_remove = []string{}

		now = time.Now()
		if _, ok := active_members[msg.AuthorID]; ok {
			temp := active_members[msg.AuthorID]
			temp.last_time = msg.Date
			temp.messages++
			active_members[msg.AuthorID] = temp
		} else {
			active_members[msg.AuthorID] = ActiveMembersStruct{
				start_time: msg.Date,
				last_time:  msg.Date,
				messages:   1,
			}
		}
		consecutiveTimeTimer1 += time.Since(now)

		now = time.Now()
		for userID, info := range active_members {
			timeDifference := msg.Date.Sub(info.last_time)

			if timeDifference.Seconds()/60 > 10 {
				timeDelta := info.last_time.Sub(info.start_time)

				if timeDelta.Seconds() == 0 {
					timeDelta = time.Duration(time.Minute * 2)
				}
				if info.messages*10 < int(timeDelta.Seconds()/60) {
					timeDelta = time.Duration(int64(time.Minute) * 2 * int64(info.messages))
				}

				ccB.consecutiveTime[userID] = append(ccB.consecutiveTime[userID], timeDelta.Seconds())

				to_remove = append(to_remove, userID)
			}
		}
		consecutiveTimeTimer2 += time.Since(now)

		now = time.Now()
		for _, userID := range to_remove {
			delete(active_members, userID)
		}
		consecutiveTimeTimer3 += time.Since(now)

		// roleDistribution calculations
		var to_add []string
		var contains_java bool = false
		var contains_bedrock bool = false

		now = time.Now()
		for _, role := range msg.RolesIDs {
			if _, ok := roles[role]; ok {
				to_add = append(to_add, role)
				ccB.roleDistribution[role] = ccB.roleDistribution[role] + 1
			}
			if role == javaRole {
				contains_java = true
			}
			if role == bedrockRole {
				contains_bedrock = true
			}
		}
		roleDistributionTimer1 += time.Since(now)

		now = time.Now()
		if contains_java && !contains_bedrock {
			ccB.roleDistribution["Java Only"] = ccB.roleDistribution["Java Only"] + 1
		} else if !contains_java && contains_bedrock {
			ccB.roleDistribution["Bedrock Only"] = ccB.roleDistribution["Bedrock Only"] + 1
		} else if contains_java && contains_bedrock {
			ccB.roleDistribution["Dual"] = ccB.roleDistribution["Dual"] + 1
		}

		if len(to_add) == 0 {
			ccB.roleDistribution["No Roles"] = ccB.roleDistribution["No Roles"] + 1
		}
		roleDistributionTimer2 += time.Since(now)

		// hourlyActivity
		ccB.hourlyActivity[msg.Date.Hour()] = ccB.hourlyActivity[msg.Date.Hour()] + 1
	}

	logging.LOGGING(fmt.Sprintf("Counting total messages per person took: %v", totalMessagesTimer), "CCI.processDB")
	logging.LOGGING(fmt.Sprintf("Adding or updating consecutive struct took: %v", consecutiveTimeTimer1), "CCI.processDB")
	logging.LOGGING(fmt.Sprintf("Looping over active members took: %v", consecutiveTimeTimer2), "CCI.processDB")
	logging.LOGGING(fmt.Sprintf("Removing inactive members took: %v", consecutiveTimeTimer3), "CCI.processDB")
	logging.LOGGING(fmt.Sprintf("Looping over user roles took: %v", roleDistributionTimer1), "CCI.processDB")
	logging.LOGGING(fmt.Sprintf("Checking special roles took: %v", roleDistributionTimer1), "CCI.processDB")


}

func createData(ccR cubeCounterRequest) *cubeCounterData {
	usernamesDB, e := sql.Open("sqlite3", config.CC.CubePath+"usernames.sql")
	if e != nil {
		logging.ERROR("An error occurred trying to open the username database."+e.Error(), "CubeCounter.createImg")
		return &cubeCounterData{}
	}
	defer usernamesDB.Close()

	data, e := usernamesDB.Query("SELECT user_id, name FROM usernames;")
	if e != nil {
		logging.ERROR("An error occurred trying to prepare the username database."+e.Error(), "CubeCounter.createImg")
		return &cubeCounterData{}
	}

	usernames := make(map[string]string)
	for data.Next() {
		var username string
		var user_id string
		data.Scan(&user_id, &username)
		usernames[user_id] = username
	}

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

	for _, c := range ccR.channelIDs {
		processDB(c, &ccB, ccR, usernames)
	}

	return &ccB
}
