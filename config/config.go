package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config конечная структура, которя содержит
// информацию для подключения к базе, какие номера
// надо мониторить, как отправить sms и кому отправит
type Config struct {
	Token    string `json:"token"`
	Dbconfig string `json:"db_config"`
	Dids     []Did  `json:"did_numbers"`
	Smsurl   string `json:"smsurl"`
	Proxy    string `json:"proxy"`
	Users    string `json:"users_file"`
	Sms      string `json:"sms_template"`
	Period   int64  `json:"period_sec"`
}

// Did триггер, при вызове этого номера отправляются
// уведомления указанным пользователям
type Did struct {
	Number string   `json:"number"`
	Users  []string `json:"users"`
}

// GetConfig открывает указанный файл и загружает конфиг
// удобном для себя формате
func GetConfig(configFile string) (*Config, error) {
	log.Println("GetConfig", configFile)
	result := new(Config)
	conf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err.Error())
		return result, err
	}
	err = json.Unmarshal(conf, result)
	if err != nil {
		log.Println(err.Error())
		return result, err
	}
	log.Printf("%#v\n", result)
	return result, nil
}
