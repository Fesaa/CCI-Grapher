package main

import (
	"os"
	"os/signal"
	"syscall"

	"cci_grapher/config"
	"cci_grapher/cubecounterimage"
	"cci_grapher/db"
	"cci_grapher/logging"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	bot, e := discordgo.New("Bot " + config.Discord.Token)
	if e != nil {
		logging.FATAL("Error creating Discord session: "+e.Error(), "main.run")
	}
	bot.Identify.Intents = discordgo.IntentsAll
	bot.AddHandler(cubecounterimage.CCI)

	e = bot.Open()
	if e != nil {
		logging.FATAL("Error opening Discord session: "+e.Error(), "main.run")
	}
	logging.INFO(bot.State.User.Username+" has started up with session ID "+bot.State.SessionID, "main.run")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	logging.INFO("Shutting down", "main.run")
	db.Disconnect()
	err := bot.Close()
	if err != nil {
		logging.ERROR("Bot did not close correctly", "main.run")
		return
	}
}

func main() {
	config.LoadConfig("./config.json")
	logging.SetUpLogging(logging.LoggingLevel(config.Logging))
	if !db.Connect() {
		logging.FATAL("Could not connect to database", "main.main")
	}
	db.Init()
	Run()
}
