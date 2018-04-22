package main

import (
	"log"
	"os"

	"github.com/Ravillatypov/missedcall/asterisk"
	"github.com/Ravillatypov/missedcall/config"
	"github.com/Ravillatypov/missedcall/notification"
)

func main() {
	cfg, err := config.GetConfig(os.Args[1])
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	missedcalls := asterisk.Load(cfg.Dbconfig)
	notify, err := notification.Init(cfg.Token, cfg.Proxy, "звонок от %s", &cfg.Smsurl)
	if err != nil {
		log.Panic(err.Error())
	}
	notify.SendSMS(missedcalls, cfg.Dids)
	notify.SendTG(missedcalls, cfg.Dids)
}
