package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigDiscord struct {
	Token         string `json:"token"`
	ApplicationID string `json:"application_id"`
	Prefix        string `json:"prefix"`
}
type ConfigCC struct {
	CubePath   string   `json:"cube_path"`
	ChannelIDs []string `json:"channel_ids"`
}

type Config struct {
	Discord ConfigDiscord `json:"discord"`
	CC      ConfigCC      `json:"cc"`
	Logging int           `json:"logging"`
}

var Discord ConfigDiscord
var CC ConfigCC
var Logging int

func LoadConfig(path string) {
	file, e := os.ReadFile(path)
	if e != nil {
		log.Panic(e)
	}

	var c Config

	e = json.Unmarshal(file, &c)
	if e != nil {
		log.Panic(e)
	}

	Discord = c.Discord
	CC = c.CC
	Logging = c.Logging
}
