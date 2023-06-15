package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

type Server struct {
	RedirectURL string `json:"redirect_url"`
	Port        int    `json:"port"`
}

type Mysql struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	DBName   string `json:"db_name"`
}

type Feishu struct {
	AppID        string `json:"app_id"`
	AppSecret    string `json:"app_secret"`
	EncryptKey   string `json:"encrypt_key"`
	Verification string `json:"verification"`
}

type Config struct {
	Feishu Feishu `json:"feishu"`
	Server Server `json:"server"`
	Mysql  Mysql  `json:"mysql"`
}

func (m *Mysql) Dns() string {
	return m.UserName + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(int(m.Port)) + ")/" + m.DBName
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
