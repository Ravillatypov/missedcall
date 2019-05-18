package notification

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"fmt"

	"strings"

	"github.com/Ravillatypov/missedcall/asterisk"
	"github.com/Ravillatypov/missedcall/config"
	"github.com/Ravillatypov/missedcall/userlist"
	"gopkg.in/telegram-bot-api.v4"
)

type Notify struct {
	client    *http.Client
	Bot       *tgbotapi.BotAPI
	smsurl    string
	sms       string
	viocemail string
}

func Init(token, proxy, sms, smsurl, viocemail string) (*Notify, error) {
	httpProxy, err := url.Parse(proxy)
	if err != nil {
		log.Println(err.Error())
		return &Notify{}, err
	}
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(httpProxy)}}
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		log.Println(err.Error())
		result := &Notify{client: httpClient, smsurl: smsurl, sms: sms, viocemail: viocemail}
		log.Println(result)
		return result, err
	}
	result := &Notify{client: httpClient, Bot: bot, smsurl: smsurl, sms: sms, viocemail: viocemail}
	log.Println(result)
	return result, nil
}

func (n *Notify) SendSMS(calls []asterisk.Missed, dids []config.Did, users []userlist.User) {
	if len(calls) == 0 || n.smsurl == "" {
		return
	}
	log.Printf("calls: %#v\ndids: %#v\nusers: %#v\n", calls, dids, users)
	for _, call := range calls {
		for _, did := range dids {
			if call.Did == did.Number {
				for _, user := range users {
					if contain(did.Users, user.Name) && len(user.Phone) == 11 {
						msg := fmt.Sprintf(n.sms, call.Src)
						request := strings.Replace(n.smsurl, "__PHONE__", user.Phone, -1)
						request = strings.Replace(request, "__MESSAGE__", msg, -1)
						request = strings.Replace(request, "+", "%2B", -1)
						request = strings.Replace(request, " ", "+", -1)
						resp, err := http.Get(request)
						if err != nil {
							log.Println(err.Error())
						} else {
							bytes := make([]byte, 0)
							resp.Body.Read(bytes)
							log.Println(string(bytes))
							resp.Body.Close()
						}
						time.Sleep(time.Duration(1000000000))
					}
				}
			}
		}
	}

}

func (n *Notify) SendTG(calls []asterisk.Missed, dids []config.Did, users []userlist.User) {
	if len(calls) == 0 || n.Bot == nil {
		return
	}
	for _, did := range dids {
		for _, user := range getUsersByDid(users, did) {
			numbers := []string{}
			for _, call := range getCallsByDid(calls, did) {
				if call.IsVoiceMail() {
					msg := tgbotapi.NewMessage(user.Tgid, fmt.Sprintf(n.viocemail, call.Src))
					send, _ := n.Bot.Send(msg)
					voice := tgbotapi.NewVoiceUpload(user.Tgid, call.GetFilePath())
					voice.ReplyToMessageID = send.MessageID
					n.Bot.Send(voice)
				} else {
					if !contain(numbers, call.Src) {
						numbers = append(numbers, "+"+call.Src)
					}
				}
			}
			if len(numbers) > 0 {
				nums := "\n" + strings.Join(numbers, "\n")
				msg := tgbotapi.NewMessage(user.Tgid, fmt.Sprintf(n.sms, nums))
				n.Bot.Send(msg)
			}
		}
	}
}

//func (n *Notify) SendTG(calls []asterisk.Missed, dids []config.Did) error {}
func contain(slice []string, item string) bool {
	for _, it := range slice {
		if it == item {
			return true
		}
	}
	return false
}

func (n *Notify) Updates(ulist *userlist.UserList) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates, err := n.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err.Error())
	}
	for update := range updates {
		if update.Message != nil {
			ulist.SetChatID(update.Message.Chat.UserName, update.Message.Chat.ID)
		}
	}
}

func getCallsByDid(calls []asterisk.Missed, did config.Did) (filtered []asterisk.Missed) {
	for _, call := range calls {
		if call.Did == did.Number {
			filtered = append(filtered, call)
		}
	}
	return
}

func getUsersByDid(users []userlist.User, did config.Did) (filtered []userlist.User) {
	for _, user := range users {
		if contain(did.Users, user.Name) && user.Tgid != 0 {
			filtered = append(filtered, user)
		}
	}
	return
}
