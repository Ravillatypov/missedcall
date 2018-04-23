package userlist

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// User Информация о получателе уведомления
type User struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Tgusername string `json:"tg_username"`
	Tgid       int64  `json:"tg_id"`
}

// UserList список всех пользователей
type UserList struct {
	Fname string `json:"filename"`
	List  []User `json:"users"`
}

// LoadUsers загружает пользователей из файла
func LoadUsers(usersFile string) (*UserList, error) {
	log.Println("Load Users", usersFile)
	userlist := &UserList{}
	file, err := ioutil.ReadFile(usersFile)
	if err != nil {
		log.Println(err.Error())
		return userlist, err
	}
	err = json.Unmarshal(file, userlist)
	if err != nil {
		log.Println(err.Error())
		return userlist, err
	}
	log.Printf("%#v\n", userlist)
	return userlist, nil
}

// UserByName найти пользователя по имени
func (u *UserList) UserByName(name string) User {
	log.Printf("UserByName: %s", name)
	for _, user := range u.List {
		if user.Name == name {
			log.Println(user)
			return user
		}
	}
	return User{}
}

// SetChatID установить chat_id по tgname
func (u *UserList) SetChatID(tgname string, chatid int64) {
	log.Printf("SetChatID: name=%s\tchat_id=%d\n", tgname, chatid)
	for _, user := range u.List {
		if user.Tgusername == tgname {
			user.Tgid = chatid
			log.Printf("%#v\n", user)
		}
	}
}

// Save сохраняет пользователей в файл
func (u *UserList) Save() error {
	log.Println("Save users", u.Fname)
	bytes, err := json.Marshal(u)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	fileStatus, err := os.Stat(u.Fname)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = ioutil.WriteFile(u.Fname, bytes, fileStatus.Mode())
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
