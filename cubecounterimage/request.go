package cubecounterimage

import (
	"cci_grapher/config"
	"cci_grapher/db"
	"cci_grapher/utils"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (ccR *cubeCounterRequest) handleRequest(start time.Time, db *db.DataBase) (image.Image, time.Time) {
	var ccB *cubeCounterData = ccR.createData(db)
	stop1 := time.Now()
	if ccB == nil {
		utils.ERROR("createData returned nil. Cannot proceed", "CCI.handleRequest")
		return nil, stop1
	}
	utils.LOGGING(fmt.Sprintf("createData took: %v", stop1.Sub(start)), "CCI.handleRequest")

	var imgData *imageData = ccB.toImageData()
	stop2 := time.Now()
	if imgData == nil {
		utils.ERROR("toImageData returned nil. Cannot proceed", "CCI.handleRequest")
		return nil, stop2
	}
	utils.LOGGING(fmt.Sprintf("toImageData took: %v", stop2.Sub(stop1)), "CCI.handleRequest")

	var imgArray []image.Image = imgData.getImages(len(ccR.userIDs) != 0)
	stop3 := time.Now()
	if imgArray == nil {
		utils.ERROR("toImages returned nil. Cannot proceed", "CCI.handleRequest")
		return nil, stop3
	}
	utils.LOGGING(fmt.Sprintf("toImages took: %v", stop3.Sub(stop2)), "CCI.handleRequest")

	var finalImage image.Image = ccR.getFinalImage(imgArray)
	stop4 := time.Now()
	utils.LOGGING(fmt.Sprintf("imageMerge took: %v", stop4.Sub(stop3)), "CCI.handleRequest")
	return finalImage, stop4
}

func (ccR *cubeCounterRequest) createEmbed(Author *discordgo.User, elapsed time.Duration) *discordgo.MessageEmbed {
	description := fmt.Sprintf("Start date: %v %d\nEnd Date: %v %d",
		ccR.startDate.Month().String(), ccR.startDate.Day(), ccR.endDate.Month().String(), ccR.endDate.Day())
	if len(ccR.channelIDs) != len(config.Config.ChannelIDs) {
		description += "\nChannels: " + strings.Join(ccR.channelIDs, ", ")
	}

	embed := discordgo.MessageEmbed{
		Title: "Cube Counter Request <3",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/xOWrY8G.png",
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    Author.Username,
			IconURL: Author.AvatarURL(""),
		},
		Description: description,
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
