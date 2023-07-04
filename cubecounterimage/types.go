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

type imageData struct {
	totalMessageCount     int
	totalMessagesArray    []string
	totalMessages         map[string]int
	consecutiveTimeArray  []string
	consecutiveTime       map[string]float64
	roleDistributionArray []string
	roleDistribution      map[string]int
	hourlyActivity        map[int]float64
}
