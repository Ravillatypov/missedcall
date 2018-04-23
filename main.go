package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/Ravillatypov/missedcall/asterisk"
	"github.com/Ravillatypov/missedcall/config"
	"github.com/Ravillatypov/missedcall/notification"
	"github.com/Ravillatypov/missedcall/userlist"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:\n%s config_file [seconds]", os.Args[0])
		os.Exit(-1)
	}
	cfg, err := config.GetConfig(os.Args[1])
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	sec, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Println(err.Error())
		sec = -60
	}
	missedcalls := asterisk.Load(cfg.Dbconfig, sec)
	notify, err := notification.Init(cfg.Token, cfg.Proxy, cfg.Sms, cfg.Smsurl)
	if err != nil {
		log.Println(err.Error())
	}
	userlst, err := userlist.LoadUsers(cfg.Users)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-2)
	}
	if notify.Bot != nil {
		for _, tguser := range notify.Updates() {
			userlst.SetChatID(tguser.Tgusername, tguser.Tgid)
		}
	}
	notify.SendSMS(missedcalls, cfg.Dids, userlst.List)
	notify.SendTG(missedcalls, cfg.Dids, userlst.List)
	userlst.Save()
}
