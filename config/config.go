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

type ConfigStruct struct {
	Discord    ConfigDiscord `json:"discord"`
	PsqlLink   string        `json:"psql"`
	Logging    int           `json:"logging"`
	ChannelIDs []string      `json:"channel_ids"`
}

var Config ConfigStruct
var Discord ConfigDiscord
var Logging int

func LoadConfig(path string) {
	file, e := os.ReadFile(path)
	if e != nil {
		log.Panic(e)
	}

	var c ConfigStruct

	e = json.Unmarshal(file, &c)
	if e != nil {
		log.Panic(e)
	}

	Config = c
	Discord = c.Discord
	Logging = c.Logging
}
