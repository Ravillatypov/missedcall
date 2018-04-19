package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	token    string `json:"token"`
	dbconfig string `json:"db_config"`
	dids     []Did  `json:"did_numbers"`
	smsurl   SMSUrl `json:"smsurl"`
}

type SMSUrl struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type User struct {
	name       string `json:"name"`
	phone      string `json:"phone"`
	tgusername string `json:"tg_username"`
	tgid       int64  `json:"tg_id"`
}

type Did struct {
	number string `json:"number"`
	users  []User `json:"users"`
}

// GetConfig открывает указанный файл и загружает конфиг
// удобном для себя формате
func GetConfig(configFile string) *Config {
	result := new(Config)
	conf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err.Error())
		return result
	}
	err = json.UnmarshalJSON(conf, result)
	if err != nil {
		log.Println(err.Error())
	}
	return result
}
