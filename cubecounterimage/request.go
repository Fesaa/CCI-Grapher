package cubecounterimage

import (
	"cci_grapher/logging"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"image"
	"time"
)

func handleRequest(ccR cubeCounterRequest, start time.Time) (image.Image, time.Time) {
	var ccB *cubeCounterData = createData(ccR)
	stop1 := time.Now()
	logging.LOGGING(fmt.Sprintf("createData took: %v", stop1.Sub(start)), "CCI.handleRequest")

	var imgData imageData = toImageData(ccB)
	stop2 := time.Now()
	logging.LOGGING(fmt.Sprintf("toImageData took: %v", stop2.Sub(stop1)), "CCI.handleRequest")

	var imgArray []image.Image = toImages(imgData)
	stop3 := time.Now()
	logging.LOGGING(fmt.Sprintf("toImages took: %v", stop3.Sub(stop2)), "CCI.handleRequest")
	if imgArray == nil {
		logging.ERROR("toImages returned nil. Cannot proceed", "CCI.handleRequest")
		return nil, time.Now()
	}

	var finalImage image.Image = imageMerge(imgArray, ccR)
	stop4 := time.Now()
	logging.LOGGING(fmt.Sprintf("imageMerge took: %v", stop4.Sub(stop3)), "CCI.handleRequest")
	return finalImage, stop4
}

func createEmbed(ccR cubeCounterRequest, Author *discordgo.User, elapsed time.Duration) *discordgo.MessageEmbed {
	StartDate := ccR.startDate
	EndDate := ccR.endDate
	embed := discordgo.MessageEmbed{
		Title: "Cube Counter Request",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/xOWrY8G.png",
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    Author.Username,
			IconURL: Author.AvatarURL(""),
		},
		Description: fmt.Sprintf("Start date: %v %d\nEnd Date: %v %d", StartDate.Month().String(), StartDate.Day(), EndDate.Month().String(), EndDate.Day()),
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0x6A56F6,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Creation time: %d ms", elapsed/1000000),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://cci.jpg",
		},
	}
	return &embed
}
