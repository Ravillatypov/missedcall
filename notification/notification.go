package notification

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Ravillatypov/missedcall/asterisk"
	"github.com/Ravillatypov/missedcall/config"
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
)

type Notify struct {
	client *http.Client
	bot    *tgbotapi.BotAPI
	smsurl *config.SMSUrl
	sms    string
}

func Init(token, proxy, sms string, smsurl *config.SMSUrl) (*Notify, error) {
	httpProxy, err := url.Parse(proxy)
	if err != nil {
		log.Println(err.Error())
		return &Notify{}, err
	}
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(httpProxy)}}
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		log.Println(err.Error())
		return &Notify{client: httpClient, smsurl: smsurl, sms: sms}, err
	}
	return &Notify{client: httpClient, bot: bot, smsurl: smsurl, sms: sms}, nil
}

func (n *Notify) SendSMS(calls []asterisk.Missed, dids []config.Did) error {
	for _, call := range calls {
		for _, did := range dids {
			if call.Did == did.Number {
				for _,user := range did.Users{
					if user.Phone != ""{
						msg := fmt.Sprintf(n.sms, user.Phone, call.Src)
						request := &url.URL{Path: fmt.Sprintf{n.smsurl.Url, msg)}
						n.client.Get(request.String())
					}
				}
			}
		}
	}
}

//func (n *Notify) SendTG(calls []asterisk.Missed, dids []config.Did) error {}
