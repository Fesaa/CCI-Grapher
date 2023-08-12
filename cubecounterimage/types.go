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
