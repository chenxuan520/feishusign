package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Server struct {
	RedirectURL string `json:"redirect_url"`
	Port        int    `json:"port"`
}

type Feishu struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type Config struct {
	Feishu Feishu `json:"feishu"`
	Server Server `json:"server"`
}

var (
	GlobalConfig *Config
)

func InitConfig() error {
	configFile := "config.json"
	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		data, err = ioutil.ReadFile("./config/" + configFile)
		if err != nil {
			log.Println("Read config error!")
			return err
		}
	}

	config := &Config{}

	err = json.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return err
	}

	GlobalConfig = config
	log.Println("Config " + configFile + " loaded.")
	return nil
}
