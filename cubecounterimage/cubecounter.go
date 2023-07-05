package cubecounterimage

import (
	"bytes"
	"cci_grapher/config"
	"cci_grapher/logging"
	"cci_grapher/utils"
	"fmt"
	"image/png"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const dataParse string = "2006-01-02"

func CCI(s *discordgo.Session, e *discordgo.MessageCreate) {
	content := e.Content

	if !strings.HasPrefix(content, config.Discord.Prefix+"cc") {
		return
	}

	trim := strings.TrimPrefix(content, config.Discord.Prefix+"cc")
	parts := strings.Split(strings.Trim(trim, " "), " ")

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

	var ccR cubeCounterRequest = cubeCounterRequest{
		channelIDs: channelIDs,
		startDate:  StartDate,
		endDate:    EndDate,
	}
	var start time.Time = time.Now()
	finalImage, stop4 := handleRequest(ccR, start)

	var b bytes.Buffer
	if err := png.Encode(&b, finalImage); err != nil {
		logging.ERROR("An error occurred trying to convert the image to a reader:\n"+err.Error(), "cc._cci")
		return
	}
	stop5 := time.Now()
	logging.LOGGING(fmt.Sprintf("Encoding took: %v", stop5.Sub(stop4)), "CCI")
	var elapsed time.Duration = time.Since(start)

	_, err := s.ChannelMessageSendComplex(e.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:        "cci.jpg",
				ContentType: "image/jpg",
				Reader:      &b,
			},
		},
		Embeds: []*discordgo.MessageEmbed{createEmbed(ccR, e.Author, elapsed)},
	})
	if err != nil {
		logging.ERROR(fmt.Sprintf("Could not send message in %s", e.ChannelID), "cubecounter.CCI")
		return
	}

}
