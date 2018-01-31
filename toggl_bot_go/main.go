package main

import (
	"./togglbot"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
)

const (
	ACCEPT_USER = "YOUR_SLACK_ACCOUNT_EMAIL"
)

func run(api *slack.Client) int {
	rtm := api.NewRTM()
	togglInfo := &togglbot.TogglManager{}
	if togglInfo.Init() {
		os.Exit(0)
	}
	go rtm.ManageConnection()
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("Start")

			case *slack.MessageEvent:
				identity := ev.User
				user, err := api.GetUserInfo(identity)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
				if user.Profile.Email == ACCEPT_USER {
					reaction := togglbot.MessageCheck(ev.Text)
					var responseMsg = "失敗したぜ"
					if reaction != "NoMatch" {
						responseMsg = togglbot.MessageEvent(reaction, ev.Text, togglInfo)
						rtm.SendMessage(rtm.NewOutgoingMessage(responseMsg, ev.Channel))
					}
					//rtm.SendMessage(rtm.NewOutgoingMessage(responseMsg, ev.Channel))
				}
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		}
	}
}

func main() {
	api := slack.New("YOUR_SLACK_API_TOKEN")
	os.Exit(run(api))
}
