package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"github.com/Ravillatypov/missedcall/asterisk"
	"github.com/Ravillatypov/missedcall/config"
	"github.com/Ravillatypov/missedcall/notification"
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
	notify, err := notification.Init(cfg.Token, cfg.Proxy, "звонок от %s", &cfg.Smsurl)
	if err != nil {
		log.Println(err.Error())
	}
	notify.SendSMS(missedcalls, cfg.Dids)
	notify.SendTG(missedcalls, cfg.Dids)
}
