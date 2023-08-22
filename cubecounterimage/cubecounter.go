package cubecounterimage

import (
	"bytes"
	"cci_grapher/db"
	"cci_grapher/config"
	"cci_grapher/utils"
	"fmt"
	"image/png"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const dataParse string = "2006-01-02"

func CCI(s *discordgo.Session, e *discordgo.MessageCreate, db *db.DataBase) {
	content := e.Content

	if !strings.HasPrefix(content, config.Discord.Prefix+"cc") {
		return
	}

	trim := strings.TrimPrefix(content, config.Discord.Prefix+"cc")
	parts := strings.Split(strings.Trim(trim, " "), " ")

	now := time.Now()

	var channelIDs = config.Config.ChannelIDs
    var userIDS = []string{}
	var defaultChannelIDs bool = true
	var StartDate time.Time = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(time.Hour * 24 * -7)
	var EndDate time.Time = time.Now()
	var index int = 0

	for i, p := range parts {
		if utils.InStringArray(config.Config.ChannelIDs, p) {
			if defaultChannelIDs {
				channelIDs = []string{}
				defaultChannelIDs = false
			}
			channelIDs = append(channelIDs, p)
		} else {
            b, e := utils.MaybeDiscordID(p)
            if b == true && e == nil {
                userIDS = append(userIDS, p)
            } else {
			    index = i
			    break
            }
		}
	}

    if len(userIDS) > 25 {
        return
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
        userIDs:    userIDS,
	}
	var start time.Time = time.Now()
	finalImage, stop4 := ccR.handleRequest(start, db)
	if finalImage == nil {
		utils.ERROR("An error occurred trying to handle the request, returned nil", "cc._cci")
		return
	}

	var b bytes.Buffer
	if err := png.Encode(&b, finalImage); err != nil {
		utils.ERROR("An error occurred trying to convert the image to a reader:\n"+err.Error(), "cc._cci")
		return
	}
	stop5 := time.Now()
	utils.LOGGING(fmt.Sprintf("Encoding took: %v", stop5.Sub(stop4)), "CCI")
	var elapsed time.Duration = time.Since(start)

	_, err := s.ChannelMessageSendComplex(e.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:        "cci.jpg",
				ContentType: "image/jpg",
				Reader:      &b,
			},
		},
		Embed: ccR.createEmbed(e.Author, elapsed),
	})
	if err != nil {
		utils.ERROR(fmt.Sprintf("Could not send message in %s", e.ChannelID), "cubecounter.CCI")
	}
}
