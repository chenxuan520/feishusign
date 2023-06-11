package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Feishu struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type Config struct {
	Feishu Feishu `json:"feishu"`
}

var (
	GlobalConfig *Config
)

func init() {
	configFile := "config.json"
	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		data, err = ioutil.ReadFile("./config/" + configFile)
		if err != nil {
			log.Println("Read config error!")
			log.Panic(err)
			return
		}
	}

	config := &Config{}

	err = json.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return
	}

	GlobalConfig = config
	log.Println("Config " + configFile + " loaded.")
}
