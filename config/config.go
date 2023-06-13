package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)


type ConfigDiscord struct {
	Token         string               `json:"token"`
	ApplicationID string               `json:"application_id"`
}
type ConfigCC struct {
	CubePath   string   `json:"cube_path"`
	ChannelIDs []string `json:"channel_ids"`
}

type Config struct {
	Discord     ConfigDiscord     `json:"discord"`
	CC          ConfigCC          `json:"cc"`
}

var Discord ConfigDiscord
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
	CC = c.CC
}
