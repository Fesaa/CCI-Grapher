package cubecounterimage

import "time"

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

type HourMinute struct {
	Hour   int
	Minute int
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

func (ccB *cubeCounterData) AddRowInfo(msg MessageEntry, activeMembers map[string]ActiveMembersStruct) {
	ccB.totalMessageCount = ccB.totalMessageCount + 1

	// Adding to totalMessages
	if _, ok := ccB.totalMessages[msg.AuthorID]; ok {
		ccB.totalMessages[msg.AuthorID] = ccB.totalMessages[msg.AuthorID] + 1
	} else {
		ccB.totalMessages[msg.AuthorID] = 1
	}

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

			ccB.consecutiveTime[userID] = append(ccB.consecutiveTime[userID], timeDelta.Seconds())

			toRemove = append(toRemove, userID)
		}
	}

	for _, userID := range toRemove {
		delete(activeMembers, userID)
	}

	// roleDistribution calculations
	var toAdd []string
	var containsJava bool = false
	var containsBedrock bool = false

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

	// hourlyActivity
	ccB.hourlyActivity[msg.Date.Hour()] = ccB.hourlyActivity[msg.Date.Hour()] + 1
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
