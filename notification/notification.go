package notification

import (
	"log"
	"net/http"
	"net/url"

	"fmt"

	"strings"

	"github.com/Ravillatypov/missedcall/asterisk"
	"github.com/Ravillatypov/missedcall/config"
	"gopkg.in/telegram-bot-api.v4"
)

type Notify struct {
	client *http.Client
	bot    *tgbotapi.BotAPI
	smsurl *config.SMSUrl
	sms    string
}

func Init(token, proxy, sms string, smsurl *config.SMSUrl) (*Notify, error) {
	log.Println("Init notification")
	log.Println(sms, proxy)
	httpProxy, err := url.Parse(proxy)
	if err != nil {
		log.Println(err.Error())
		return &Notify{}, err
	}
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(httpProxy)}}
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		log.Println(err.Error())
		result := &Notify{client: httpClient, smsurl: smsurl, sms: sms}
		log.Println(result)
		return result, err
	}
	result := &Notify{client: httpClient, bot: bot, smsurl: smsurl, sms: sms}
	log.Println(result)
	return result, nil
}

func (n *Notify) SendSMS(calls []asterisk.Missed, dids []config.Did) {
	log.Println("SendSMS")
	if len(calls) == 0 || n.smsurl != nil {
		return
	}
	if n.smsurl.Type == "GET" {
		for _, call := range calls {
			for _, did := range dids {
				if call.Did == did.Number {
					for _, user := range did.Users {
						if len(user.Phone) == 11 {
							msg := fmt.Sprintf(n.sms, call.Src)
							log.Println(msg)
							request := fmt.Sprintf(n.smsurl.Url, user.Phone, msg)
							request = strings.Replace(request, " ", "+", -1)
							log.Println(request)
							resp, err := http.Get(request)
							if err != nil {
								log.Println(err.Error())
							} else {
								resp.Body.Close()
							}
						}
					}
				}
			}
		}
	} //else {
	// 	for _, call := range calls {
	// 		for _, did := range dids {
	// 			if call.Did == did.Number {
	// 				for _, user := range did.Users {
	// 					if len(user.Phone) == 11 {
	// 						msg := fmt.Sprintf(n.sms, call.Src)
	// 						log.Println(msg)
	// 						request := &url.URL{Path: fmt.Sprintf(n.smsurl.Url, user.Phone, msg)}
	// 						log.Println(request.String())
	// 						http.Post(request.String())
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func (n *Notify) SendTG(calls []asterisk.Missed, dids []config.Did) {
	log.Println("SendTG")
	if len(calls) == 0 || n.bot == nil {
		return
	}
	for _, call := range calls {
		for _, did := range dids {
			if call.Did == did.Number {
				for _, user := range did.Users {
					if user.Tgid != 0 {
						msg := tgbotapi.NewMessage(user.Tgid, fmt.Sprintf(n.sms, call.Src))
						log.Println(msg)
						n.bot.Send(msg)
					}
				}
			}
		}
	}
}

//func (n *Notify) SendTG(calls []asterisk.Missed, dids []config.Did) error {}
