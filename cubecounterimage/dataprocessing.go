package cubecounterimage

import (
	"cci_grapher/utils"
)

func toImageData(ccB []*cubeCounterData) *imageData {
	var totalMessageCount int = 0
	var totalMessages map[string]int = make(map[string]int)
	var consecutiveTime map[string]float64 = make(map[string]float64)
	var roleDistribution map[string]int = make(map[string]int)
	var hourlyActivity map[interface{}]float64 = make(map[interface{}]float64)

	for _, cc := range ccB {
		if cc == nil {
			continue
		}
		totalMessageCount += cc.totalMessageCount
		for k, v := range cc.totalMessages {
			if _, ok := totalMessages[k]; ok {
				
			} else {
				totalMessages[k] = v
			}
		}
		for k, v := range cc.consecutiveTime {
			consecutiveTime[k] += utils.SumOfFloat64Array(v) / 3600
		}
		for k, v := range cc.roleDistribution {
			roleDistribution[k] += v
		}
		for k, v := range cc.hourlyActivity {
			hourlyActivity[k] += float64(v)
		}
	}

	var c map[string]float64 = make(map[string]float64)
	for k, v := range consecutiveTime {
		c[k] = v / 3600
	}

	var r = make(map[string]int)
	for k, v := range roleDistribution {
		r[roles[k]] = int(float64(v) / float64(totalMessageCount) * 100)
	}

	var h = make(map[interface{}]float64)
	for k, v := range hourlyActivity {
		h[k] = v / float64(totalMessageCount) * 100
	}

	return &imageData{
		totalMessageCount:     totalMessageCount,
		totalMessagesArray:    utils.TopNOfIntMap(totalMessages, 25),
		totalMessages:         totalMessages,
		consecutiveTimeArray:  utils.TopNOfFloat64Map(c, 25),
		consecutiveTime:       c,
		roleDistributionArray: utils.TopNOfIntMap(r, 25),
		roleDistribution:      r,
		hourlyActivity:        h,
	}

}
