package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

var TestMode = false

type Server struct {
	SignRedirectURL  string `json:"sign_redirect_url"`
	LoginRedirectURL string `json:"login_redirect_url"`
	Port             int    `json:"port"`
	StaticPath       string `json:"static_path"`
}

type Mysql struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	DBName   string `json:"db_name"`
}

type Feishu struct {
	AppID        string   `json:"app_id"`
	AppSecret    string   `json:"app_secret"`
	EncryptKey   string   `json:"encrypt_key"`
	Verification string   `json:"verification"`
	Root         []string `json:"root"`
}

type Sign struct {
	HashSalt       string        `json:"hash_salt"`
	ChangeTime     time.Duration `json:"change_time"`
	ExpireDuration time.Duration `json:"expire_duration"`
	JwtToken       string        `json:"jwt_token"`
	Issuer         string        `json:"issuer"`
	FolderToken    string        `json:"folder_token"`
}

type Config struct {
	Feishu Feishu `json:"feishu"`
	Server Server `json:"server"`
	Mysql  Mysql  `json:"mysql"`
	Sign   Sign   `json:"sign"`
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
