package cubecounterimage

import (
	"cci_grapher/utils"
	"fmt"
	"image"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

var colourMap = map[int]drawing.Color{
	0: chart.ColorRed,
	1: chart.ColorOrange,
	2: chart.ColorAlternateYellow,
	3: chart.ColorAlternateGreen,
	4: chart.ColorAlternateBlue,
}

func getColour(i int, val int) drawing.Color {
	if val == 0 {
		return chart.ColorAlternateGray
	}
	return colourMap[i%5]
}

func (iD *imageData) getImages(isUsers bool) []image.Image {
	if isUsers {
		return iD.toUserImages()
	} else {
		return iD.toImages()
	}
}

func (iD *imageData) toUserImages() []image.Image {
	iD.totalMessagesArray = append(iD.totalMessagesArray, "")
	iD.totalMessages[""] = 0

	totalMessagesImage, err := iD.getTotalMessagesImage()
	if err != nil {
		utils.ERROR("Generating Total Messages Image enocuntered an error. Cannot continue.", "cubecounterimage.toImages")
		return nil
	}

	hourlyActivityImage, err := iD.getHourlyActivityImage()
	if err != nil {
		utils.ERROR("Generating Hourly Activity Image enocuntered an error. Cannot continue.", "cubecounterimage.toImages")
		return nil
	}

	var imgs []image.Image
	imgs = append(imgs, totalMessagesImage)
	imgs = append(imgs, hourlyActivityImage)
	return imgs
}

func (iD *imageData) toImages() []image.Image {
	totalMessagesImage, err := iD.getTotalMessagesImage()
	if err != nil {
		utils.ERROR("Generating Total Messages Image enocuntered an error. Cannot continue.", "cubecounterimage.toImages")
		return nil
	}

	consecutiveTimeImage, err := iD.getConsecutiveTimeImage()
	if err != nil {
		utils.ERROR("Generating Consecutive Time Image enocuntered an error. Cannot continue.", "cubecounterimage.toImages")
		return nil
	}

	roleDistributionImage, err := iD.getRoleDistributionImage()
	if err != nil {
		utils.ERROR("Generating Role Distribution Image enocuntered an error. Cannot continue.", "cubecounterimage.toImages")
		return nil
	}

	hourlyActivityImage, err := iD.getHourlyActivityImage()
	if err != nil {
		utils.ERROR("Generating Hourly Activity Image enocuntered an error. Cannot continue.", "cubecounterimage.toImages")
		return nil
	}

	var imgs []image.Image
	imgs = append(imgs, totalMessagesImage)
	imgs = append(imgs, consecutiveTimeImage)
	imgs = append(imgs, roleDistributionImage)
	imgs = append(imgs, hourlyActivityImage)
	return imgs

}

func (iD *imageData) getTotalMessagesImage() (image.Image, error) {
	var totalMessagesBars []chart.Value
	for i, v := range iD.totalMessagesArray {
		idx := i / utils.Max(len(iD.totalMessagesArray)/5, 1)
		val := iD.totalMessages[v]
		totalMessagesBars = append(totalMessagesBars, chart.Value{
			Label: v,
			Value: float64(val),
			Style: chart.Style{
				FillColor:   getColour(idx, val),
				StrokeColor: getColour(idx, val),
				DotColor:    getColour(idx, val),
			},
		})
	}
	totalMessagesChart := barChartMaker(fmt.Sprintf("Total messages (%d)", iD.totalMessageCount), totalMessagesBars)
	totalMessagesCollector := &chart.ImageWriter{}
	err := totalMessagesChart.Render(chart.PNG, totalMessagesCollector)
	if err != nil {
		utils.ERROR("Could not Render Total Messages chart:\n"+err.Error(), "cubecounterimage.getTotalMessagesImage")
		return nil, err
	}
	totalMessagesImage, err := totalMessagesCollector.Image()
	if err != nil {
		utils.ERROR("Could not collect Total Messages chart:\n"+err.Error(), "cubecounterimage.getTotalMessagesImage")
		return nil, err
	}
	return totalMessagesImage, nil
}

func (iD *imageData) getConsecutiveTimeImage() (image.Image, error) {
	var consecutiveTimeBars []chart.Value
	for i, v := range iD.consecutiveTimeArray {
		consecutiveTimeBars = append(consecutiveTimeBars, chart.Value{
			Label: v,
			Value: iD.consecutiveTime[v],
			Style: chart.Style{
				FillColor:   colourMap[(i / 5)],
				StrokeColor: colourMap[(i / 5)],
				DotColor:    colourMap[(i / 5)],
			},
		})
	}
	consecutiveTimeChart := barChartMaker("Consecutive Time (h)", consecutiveTimeBars)
	consecutiveTimeCollector := &chart.ImageWriter{}
	err := consecutiveTimeChart.Render(chart.PNG, consecutiveTimeCollector)
	if err != nil {
		utils.ERROR("Could not Render Consecutive Time chart:\n"+err.Error(), "cubecounterimage.getConsecutiveTimeImage")
		return nil, err
	}
	consecutiveTimeImage, err := consecutiveTimeCollector.Image()
	if err != nil {
		utils.ERROR("Could not collect Consecutive Time chart:\n"+err.Error(), "cubecounterimage.getConsecutiveTimeImage")
		return nil, err
	}
	return consecutiveTimeImage, nil
}

func (iD *imageData) getRoleDistributionImage() (image.Image, error) {
	var roleDistributionBars []chart.Value
	for i, v := range iD.roleDistributionArray {
		roleDistributionBars = append(roleDistributionBars, chart.Value{
			Label: v,
			Value: float64(iD.roleDistribution[v]),
			Style: chart.Style{
				FillColor:   colourMap[(i / 5)],
				StrokeColor: colourMap[(i / 5)],
				DotColor:    colourMap[(i / 5)],
			},
		})
	}
	roleDistributionChart := barChartMaker("Role Distribution (%)", roleDistributionBars)
	roleDistributionCollector := &chart.ImageWriter{}
	err := roleDistributionChart.Render(chart.PNG, roleDistributionCollector)
	if err != nil {
		utils.ERROR("Could not Render Role Distribution chart:\n"+err.Error(), "cubecounterimage.getRoleDistributionImage")
		return nil, err
	}
	roleDistributionImage, err := roleDistributionCollector.Image()
	if err != nil {
		utils.ERROR("Could not collect Role Distribution chart:\n"+err.Error(), "cubecounterimage.getRoleDistributionImage")
		return nil, err
	}
	return roleDistributionImage, nil
}

func (iD *imageData) getHourlyActivityImage() (image.Image, error) {
	var hourlyActivityBars []chart.Value
	min, max := utils.MinMaxOfMap(iD.hourlyActivity)
	for i := 0; i < 24; i++ {
		v := iD.hourlyActivity[i]
		c := 5 - int((v-min)/(max-min)*5) - 1
		if c < 0 {
			c = 0
		}
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
	hourlyActivityCollector := &chart.ImageWriter{}
	err := hourlyActivityChart.Render(chart.PNG, hourlyActivityCollector)
	if err != nil {
		utils.ERROR("Could not Render Hourly Activity chart:\n"+err.Error(), "cubecounterimage.getHourlyActivityImage")
		return nil, err
	}
	hourlyActivityImage, err := hourlyActivityCollector.Image()
	if err != nil {
		utils.ERROR("Could not collect Hourly Activity chart:\n"+err.Error(), "cubecounterimage.getHourlyActivityImage")
		return nil, err
	}
	return hourlyActivityImage, nil
}

func barChartMaker(title string, bars []chart.Value) chart.BarChart {
	return chart.BarChart{
		Title: title,
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
		Bars:     bars,
		XAxis: chart.Style{
			TextVerticalAlign:   chart.TextVerticalAlignTop,
			TextRotationDegrees: 85,
			FontSize:            13,
		},
	}
}
