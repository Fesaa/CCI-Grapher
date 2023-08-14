package cubecounterimage

import (
	"cci_grapher/utils"
	"database/sql"
	"strings"
	"time"
)

type cubeCounterRequest struct {
	channelIDs []string
	startDate  time.Time
	endDate    time.Time
}

type MessageEntry struct {
	Date     time.Time
	AuthorID string
	RolesIDs []string
}

type ActiveMembersStruct struct {
	messages  int
	startTime time.Time
	lastTime  time.Time
}

type cubeCounterData struct {
	totalMessageCount int
	totalMessages     map[string]int
	consecutiveTime   map[string][]float64
	roleDistribution  map[string]int
	hourlyActivity    map[int]int
}

func GetCubeCounterDate() cubeCounterData {
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
	return ccB
}

type CubeCounterDataType int

const (
	TotalMessageCount CubeCounterDataType = iota
	TotalMessages
	ConsecutiveTime
	RoleDistribution
	HourlyActivity
)

type Message struct {
	cubeCounterDataType CubeCounterDataType
	data                interface{}
}

func AddRowInfo(msg MessageEntry, activeMembers map[string]ActiveMembersStruct, ch chan Message) {
	ch <- Message{TotalMessageCount, nil}

	// Adding to totalMessages
	ch <- Message{TotalMessages, msg.AuthorID}

	// consecutiveTime calculations
	var toRemove []string

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

			ch <- Message{ConsecutiveTime, []interface{}{userID, timeDelta.Seconds()}}
			toRemove = append(toRemove, userID)
		}
	}

	for _, userID := range toRemove {
		delete(activeMembers, userID)
	}

	// roleDistribution calculations
	var containsJava bool = false
	var containsBedrock bool = false

    if len(msg.RolesIDs) == 1 {
        ch <- Message{RoleDistribution, "No roles"}
    } else {
        for _, role := range msg.RolesIDs {
            if _, ok := roles[role]; ok {
                ch <- Message{RoleDistribution, role}
            }
            if role == javaRole {
                containsJava = true
            }
            if role == bedrockRole {
                containsBedrock = true
            }
        }

        if containsJava && !containsBedrock {
            ch <- Message{RoleDistribution, "Java Only"}
        } else if !containsJava && containsBedrock {
            ch <- Message{RoleDistribution, "Bedrock Only"}
        } else if containsJava && containsBedrock {
            ch <- Message{RoleDistribution, "Dual"}
        }

    }


	// hourlyActivity
	ch <- Message{HourlyActivity, msg.Date.Hour()}
}

func processRows(rows *sql.Rows, userGetter map[string]string,
	activeMembers map[string]ActiveMembersStruct, ch chan Message) (string, int, error) {
	var messageID string = ""
	var rowCounter int = 0
	for rows.Next() {
		rowCounter++
		var userID string
		var time time.Time
		var rolesString string
		err := rows.Scan(&messageID, &userID, &rolesString, &time)
		if err != nil {
			utils.ERROR("An error occurred trying to scan from rows."+err.Error(), "CubeCounter.processDB")
			return "", 0, err
		}
		msg := MessageEntry{
			Date:     time,
			AuthorID: userGetter[userID],
			RolesIDs: strings.Split(rolesString, ","),
		}
		AddRowInfo(msg, activeMembers, ch)
	}
	if rowCounter < 10000 {
		messageID = ""
	}
	return messageID, rowCounter, nil
}

type imageData struct {
	totalMessageCount     int
	totalMessagesArray    []string
	totalMessages         map[string]int
	consecutiveTimeArray  []string
	consecutiveTime       map[string]float64
	roleDistributionArray []string
	roleDistribution      map[string]int
	hourlyActivity        map[interface{}]float64
}
