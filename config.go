package config

import (
	"io/ioutil"
	"log"
)

type Config struct {
	token    string `json:"token"`
	dbconfig string `json:"db_config"`
	dids     []Did  `json:"did_numbers"`
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

// UnmarshalJSON преобразовывает из json в свой тип Did
func (d *Did) UnmarshalJSON(data []byte) error {}

// UnmarshalJSON преобразовывает из json в свой тип User
func (u *User) UnmarshalJSON(data []byte) error {}

// UnmarshalJSON преобразовывает из json в свой тип Config
func (c *Config) UnmarshalJSON(data []byte) error {}

// GetConfig открывает указанный файл и загружает конфиг
// удобном для себя формате
func GetConfig(configFile string) *Config {
	result := new(Config)
	conf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err.Error())
		return result
	}
	result.UnmarshalJSON(conf)
	return result
}
