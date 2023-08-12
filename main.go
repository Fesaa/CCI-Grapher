package main

import (
	"cci_grapher/utils"
	"os"
	"os/signal"
	"syscall"

	"cci_grapher/config"
	"cci_grapher/cubecounterimage"
	"cci_grapher/db"
	"github.com/bwmarrin/discordgo"
)

func Run() {
	bot, e := discordgo.New("Bot " + config.Discord.Token)
	if e != nil {
		utils.FATAL("Error creating Discord session: "+e.Error(), "main.run")
	}
	bot.Identify.Intents = discordgo.IntentsAll

	dataBase := db.DataBase{}
	e = dataBase.Connect(config.Config.PsqlLink)
	if e != nil {
		utils.FATAL("Error connecting to database: "+e.Error(), "main.run")
	}
	e = dataBase.Init()
	if e != nil {
		utils.FATAL("Error initializing database: "+e.Error(), "main.run")
	}

	bot.AddHandler(func (s *discordgo.Session, e *discordgo.MessageCreate)  {
		go cubecounterimage.CCI(s, e, &dataBase)
	})

	e = bot.Open()
	if e != nil {
		utils.FATAL("Error opening Discord session: "+e.Error(), "main.run")
	}
	utils.SUCCESS(bot.State.User.Username+" has started up with session ID "+bot.State.SessionID, "main.run")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	utils.INFO("Shutting down", "main.run")
	err := dataBase.Disconnect()
	if err != nil {
		utils.ERROR("Database did not close correctly", "main.run")
	}
	err = bot.Close()
	if err != nil {
		utils.ERROR("Bot did not close correctly", "main.run")
		return
	}
	utils.SUCCESS("Bot closed correctly", "main.run")
}

func main() {
	config.LoadConfig("./config.json")
	utils.SetUpLogging(utils.LoggingLevel(config.Logging))
	Run()
}
