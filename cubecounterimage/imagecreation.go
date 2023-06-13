package cubecounterimage

import (
	"cci_grapher/logging"
	"cci_grapher/utils"
	"fmt"
	"image"
	"math"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func toImages(iD imageData) []image.Image {

	var colourMap = map[int]drawing.Color{
		0: chart.ColorRed,
		1: chart.ColorOrange,
		2: chart.ColorAlternateYellow,
		3: chart.ColorAlternateGreen,
		4: chart.ColorAlternateBlue,
	}

	var imgs []image.Image

	// Total messages
	var totalMessagesBars []chart.Value
	for i, v := range iD.totalMessagesArray {
		totalMessagesBars = append(totalMessagesBars, chart.Value{
			Label: v,
			Value: float64(iD.totalMessages[v]),
			Style: chart.Style{
				FillColor:   colourMap[int(i/5)],
				StrokeColor: colourMap[int(i/5)],
				DotColor:    colourMap[int(i/5)],
			},
		})
	}
	totalMessagesChart := barChartMaker(fmt.Sprintf("Total messages (%d)", iD.totalMessageCount), totalMessagesBars)
	totalMessagesCollector := &chart.ImageWriter{}
	totalMessagesChart.Render(chart.PNG, totalMessagesCollector)
	totalMessagesImage, err := totalMessagesCollector.Image()
	if err != nil {
		logging.ERROR("Could not collect Total Messages chart:\n"+err.Error(), "cubecounterimage.toImages")
	}

	// Consecutive time
	var consecutiveTimeBars []chart.Value
	for i, v := range iD.consecutiveTimeArray {
		consecutiveTimeBars = append(consecutiveTimeBars, chart.Value{
			Label: v,
			Value: float64(iD.consecutiveTime[v]),
			Style: chart.Style{
				FillColor:   colourMap[int(i/5)],
				StrokeColor: colourMap[int(i/5)],
				DotColor:    colourMap[int(i/5)],
			},
		})
	}
	consecutiveTimeChart := barChartMaker("Consecutive Time (h)", consecutiveTimeBars)
	consecutiveTimeCollector := &chart.ImageWriter{}
	consecutiveTimeChart.Render(chart.PNG, consecutiveTimeCollector)
	consecutiveTimeImage, err := consecutiveTimeCollector.Image()
	if err != nil {
		logging.ERROR("Could not collect Consecutive Time chart:\n"+err.Error(), "cubecounterimage.toImages")
	}

	// Role distribution
	var roleDistributionBars []chart.Value
	for i, v := range iD.roleDistributionArray {
		roleDistributionBars = append(roleDistributionBars, chart.Value{
			Label: v,
			Value: float64(iD.roleDistribution[v]),
			Style: chart.Style{
				FillColor:   colourMap[int(i/5)],
				StrokeColor: colourMap[int(i/5)],
				DotColor:    colourMap[int(i/5)],
			},
		})
	}
	roleDistributionChart := barChartMaker("Role Distribution (%)", roleDistributionBars)
	roleDistributionCollector := &chart.ImageWriter{}
	roleDistributionChart.Render(chart.PNG, roleDistributionCollector)
	roleDistributionImage, err := roleDistributionCollector.Image()
	if err != nil {
		logging.ERROR("Could not collect Role Distribution chart:\n"+err.Error(), "cubecounterimage.toImages")
	}

	// Hourly Activity
	var hourlyActivityBars []chart.Value
	values := make([]float64, 0, len(iD.hourlyActivity))
	for _, v := range iD.hourlyActivity {
		values = append(values, v)
	}
	max, min := utils.MaxOfFloat64Array(values), utils.MinOfFloat64Array(values)
	for i := 0; i < 24; i++ {
		v := iD.hourlyActivity[i]
		c := 5 - (int(math.Ceil((v - min) / ((max - min) / 5))))
		hourlyActivityBars = append(hourlyActivityBars, chart.Value{
			Label: fmt.Sprintf("%d", i),
			Value: v,
			Style: chart.Style{
				FillColor:   colourMap[c],
				StrokeColor: colourMap[c],
				DotColor:    colourMap[c],
			},
		})
	}
	hourlyActivityChart := barChartMaker("Hourly Activity (%)", hourlyActivityBars)
	//hourlyActivityChart.BaseValue = 0
	//hourlyActivityChart.UseBaseValue = true
	hourlyActivityCollector := &chart.ImageWriter{}
	hourlyActivityChart.Render(chart.PNG, hourlyActivityCollector)
	hourlyActivityImage, err := hourlyActivityCollector.Image()
	if err != nil {
		logging.ERROR("Could not collect Hourly Activity chart:\n"+err.Error(), "cubecounterimage.toImages")
	}

	imgs = append(imgs, totalMessagesImage)
	imgs = append(imgs, consecutiveTimeImage)
	imgs = append(imgs, roleDistributionImage)
	imgs = append(imgs, hourlyActivityImage)
	return imgs

}

func barChartMaker(t string, b []chart.Value) chart.BarChart {
	return chart.BarChart{
		Title: t,
		TitleStyle: chart.Style{
			FontColor: chart.ColorBlack,
		},
		Background: chart.Style{
			FillColor: chart.ColorAlternateGray,
			Padding: chart.Box{
				Top:    40,
				Bottom: 80,
			},
		},
		Canvas: chart.Style{
			FillColor: chart.ColorAlternateGray,
		},
		Height:   512,
		BarWidth: 30,
		Bars:     b,
		XAxis: chart.Style{
			TextVerticalAlign:   chart.TextVerticalAlignTop,
			TextRotationDegrees: 85,
			FontSize:            13,
		},
	}
}
