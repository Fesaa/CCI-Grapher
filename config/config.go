package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ConfigDiscordWebhook struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

type ConfigDiscord struct {
	Token         string               `json:"token"`
	ApplicationID string               `json:"application_id"`
	DefaultPrefix string               `json:"default_prefix"`
	Webhook       ConfigDiscordWebhook `json:"webhook"`
}

type ConfigServer struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Configcci_grapher struct {
	OwnerIds     []string `json:"owner_ids"`
	ColourString string   `json:"colour"`
	Colour       int
}

type ConfigCC struct {
	CubePath   string   `json:"cube_path"`
	ChannelIDs []string `json:"channel_ids"`
}

type Config struct {
	Discord     ConfigDiscord     `json:"discord"`
	Server      ConfigServer      `json:"server"`
	cci_grapher Configcci_grapher `json:"cci_grapher"`
	CC          ConfigCC          `json:"cc"`
}

var Discord ConfigDiscord
var Server ConfigServer
var cci_grapher Configcci_grapher
var CC ConfigCC

func LoadConfig(path string) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		log.Panic(e)
	}

	var c Config

	e = json.Unmarshal(file, &c)
	if e != nil {
		log.Panic(e)
	}

	Discord = c.Discord
	Server = c.Server
	cci_grapher = c.cci_grapher
	var i int
	fmt.Sscan(cci_grapher.ColourString, &i)
	cci_grapher.Colour = i
	CC = c.CC
}
