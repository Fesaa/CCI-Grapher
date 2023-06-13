package cubecounterimage

import (
	"bytes"
	"cci_grapher/config"
	"cci_grapher/logging"
	"cci_grapher/utils"
	"fmt"
	"image"
	"image/png"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const dataParse string = "2006-01-02"

func CCI(s *discordgo.Session, e *discordgo.MessageCreate) {
	content := e.Content

	if !strings.HasPrefix(content, "?cc") {
		return
	}

	parts := strings.Split(strings.Trim(strings.TrimPrefix(content, "?cc"), " "), " ")

	now := time.Now()

	var channelIDs = []string{"174837853778345984"}
	var defaultChannelIDs bool = true
	var StartDate time.Time = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(time.Hour * 24 * -7)
	var EndDate time.Time = time.Now()
	var index int = 0

	for i, p := range parts {
		if utils.InStringArray(config.CC.ChannelIDs, p) {
			if defaultChannelIDs {
				channelIDs = []string{}
				defaultChannelIDs = false
			}
			channelIDs = append(channelIDs, p)
		} else {
			index = i
			break
		}
	}

	if len(parts) > index {
		t, err := time.Parse(dataParse, parts[index])
		if err == nil {
			StartDate = t
		}
	}
	if len(parts) > index+1 {
		t, err := time.Parse(dataParse, parts[index+1])
		if err == nil {
			EndDate = t
		}

	}

	var start time.Time = time.Now()
	var ccR cubeCounterRequest = cubeCounterRequest{
		channelIDs: channelIDs,
		startDate:  StartDate,
		endDate:    EndDate,
	}
	var ccB *cubeCounterData = createData(ccR)
	stop1 := time.Now()
	logging.LOGGING(fmt.Sprintf("createData took: %v", stop1.Sub(start)), "CCI")

	var imgData imageData = toImageData(ccB)
	stop2 := time.Now()
	logging.LOGGING(fmt.Sprintf("toImageData took: %v", stop2.Sub(stop1)), "CCI")

	var imgArray []image.Image = toImages(imgData)
	stop3 := time.Now()
	logging.LOGGING(fmt.Sprintf("toImages took: %v", stop3.Sub(stop2)), "CCI")
	if imgArray == nil {
		logging.ERROR("toImages returned nil. Cannot proceed", "cc._cci")
		return
	}

	var finalImage image.Image = imageMerge(imgArray, ccR)
	stop4 := time.Now()
	logging.LOGGING(fmt.Sprintf("imageMerge took: %v", stop4.Sub(stop3)), "CCI")

	var b bytes.Buffer
	if err := png.Encode(&b, finalImage); err != nil {
		logging.ERROR("An error occurred trying to convert the image to a reader:\n"+err.Error(), "cc._cci")
		return
	}
	stop5 := time.Now()
	logging.LOGGING(fmt.Sprintf("Encoding took: %v", stop5.Sub(stop4)), "CCI")
	var elapsed time.Duration = time.Since(start)

	s.ChannelMessageSendComplex(e.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:        "cci.jpg",
				ContentType: "image/jpg",
				Reader:      &b,
			},
		},
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Cube Counter Request",
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://i.imgur.com/xOWrY8G.png",
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    e.Author.Username,
					IconURL: e.Author.AvatarURL(""),
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
			},
		},
	})

}
