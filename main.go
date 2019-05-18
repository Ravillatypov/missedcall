package main

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/Ravillatypov/missedcall/asterisk"
	"github.com/Ravillatypov/missedcall/config"
	"github.com/Ravillatypov/missedcall/notification"
	"github.com/Ravillatypov/missedcall/userlist"
)

func main() {
	log.Println("Start")
	if len(os.Args) < 2 {
		fmt.Printf("Usage:\n%s <config_file>", os.Args[0])
		os.Exit(-1)
	}
	cfg, err := config.GetConfig(os.Args[1])
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	missedcalls := asterisk.Load(cfg.Dbconfig, cfg.Period)
	notify, err := notification.Init(cfg.Token, cfg.Proxy, cfg.Sms, cfg.Smsurl, cfg.Voice)
	if err != nil {
		log.Println(err.Error())
	}
	userlst, err := userlist.LoadUsers(cfg.Users)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-2)
	}
	if notify.Bot != nil {
		go notify.Updates(userlst)
	}
	time.Sleep(time.Duration(20000000000))
	notify.SendSMS(missedcalls, cfg.Dids, userlst.List)
	notify.SendTG(missedcalls, cfg.Dids, userlst.List)
	userlst.Save()
}
