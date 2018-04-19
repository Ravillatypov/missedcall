package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"os"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	today := time.Now()
	date := ""
	logfilename := "C:/Program Files (x86)/Cobian Backup 11/Logs/log "
	if today.Month() < 10 {
		date = fmt.Sprintf("%d-0%d-%d", today.Year(), today.Month(), today.Day())
	} else {
		date = fmt.Sprintf("%d-%d-%d", today.Year(), today.Month(), today.Day())
	}
	logfilename += date + ".txt"
	fmt.Println(logfilename)
	fmt.Println(strings.Join(os.Args[1:], " "))
	bytes, err := ioutil.ReadFile(logfilename)
	if err != nil {
		fmt.Println("can't open file")
		return
	}
	errors := ""
	logfile := strings.Split(string(bytes), "\n")
	for _, line := range logfile {
		if strings.HasPrefix(line, "ERR") {
			errors += "\n" + line
		}
	}
	errors = strings.Trim(errors, date)
	errors = strings.Trim(errors, "ERR")
	proxyUrl, err := url.Parse("socks5://suz.iqvision.pro:9988")
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	bot, err := tgbotapi.NewBotAPIWithClient("382966119:AAGtrvgdGNtjPo0C4gCdsXwqLBO-98QXtsI", myClient)
	if err != nil {
		fmt.Println("can't authorize")
	}
	msg := tgbotapi.NewMessage(-221172754, strings.Join(os.Args[2:], " "))
	if len(errors) > 0 {
		msg.Text += errors
	} else {
		msg.Text += "Бекап сделан без ошибок"
	}
	bot.Send(msg)
}
