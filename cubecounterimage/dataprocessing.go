package cubecounterimage

import (
	"cci_grapher/utils"
)

func toImageData(ccB *cubeCounterData) imageData {

	var c map[string]float64 = make(map[string]float64)
	for k, v := range ccB.consecutiveTime {
		c[k] = utils.SumOfFloat64Array(v) / 3600
	}

	var r = make(map[string]int)
	for k, v := range ccB.roleDistribution {
		r[roles[k]] = int(float64(v) / float64(ccB.totalMessageCount) * 100)
	}

	var h = make(map[interface{}]float64)
	for k, v := range ccB.hourlyActivity {
		h[k] = float64(v) / float64(ccB.totalMessageCount) * 100
	}

	return imageData{
		totalMessageCount:     ccB.totalMessageCount,
		totalMessagesArray:    utils.TopNOfIntMap(ccB.totalMessages, 25),
		totalMessages:         ccB.totalMessages,
		consecutiveTimeArray:  utils.TopNOfFloat64Map(c, 25),
		consecutiveTime:       c,
		roleDistributionArray: utils.TopNOfIntMap(r, 25),
		roleDistribution:      r,
		hourlyActivity:        h,
	}

}
