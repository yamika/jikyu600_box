package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"./toggl_bot"
	"log"
	"os"
)

const (
	acceptuser = "YOUR_SLACK_ACCOUNT_EMAIL"
)
func run(api *slack.Client) int {
	rtm := api.NewRTM()
	toggl_info := &toggl_bot.Toggl_manager{}
	if toggl_info.Init(){
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
				if user.Profile.Email == acceptuser{
					reaction := toggl_bot.Message_check(ev.Text)
					var response_msg = "失敗したぜ"
					if reaction != "NoMatch" {
						response_msg = toggl_bot.MessageEvent(reaction,ev.Text,toggl_info)
						rtm.SendMessage(rtm.NewOutgoingMessage(response_msg, ev.Channel))
					}
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
