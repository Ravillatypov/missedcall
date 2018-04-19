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
	token    string `json:"token"`
	dbconfig string `json:"db_config"`
	dids     []Did  `json:"did_numbers"`
	smsurl   SMSUrl `json:"smsurl"`
	proxy    string `json:"proxy"`
}

// SMSUrl содержит url и каким методом надо вызывать (GET, POST)
type SMSUrl struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

// User Информация о получателе уведомления
type User struct {
	name       string `json:"name"`
	phone      string `json:"phone"`
	tgusername string `json:"tg_username"`
	tgid       int64  `json:"tg_id"`
}

// Did триггер, при вызове этого номера отправляются
// уведомления указанным пользователям
type Did struct {
	number string `json:"number"`
	users  []User `json:"users"`
}

// GetConfig открывает указанный файл и загружает конфиг
// удобном для себя формате
func GetConfig(configFile string) (*Config, error) {
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
	return result, nil
}
