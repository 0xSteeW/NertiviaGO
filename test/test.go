package main

import (
	"fmt"
	"log"
	nertivia "nertivia/api"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	client := nertivia.New("BOT_TOKEN",5)
	err := client.Open()
	if err != nil {
		log.Fatal(err)
		return
	}
	client.OnMessage(NewMessage)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit
	client.Close()
}

func NewMessage(session *nertivia.Session, messageCreate *nertivia.MessageCreate) {
	if len(messageCreate.Message.Mentions) != 0 {
		if messageCreate.Message.Mentions[0].ID == "6681535949026889728" {
			session.ChannelMessageSend(messageCreate.Message.ChannelID, "Hello! I'm Gomez. You can use my prefix go!")
			return
		}
	}
	router := nertivia.NewRouter("go!")
	router.Add("gomez", func() {
		session.ChannelMessageSend(messageCreate.Message.ChannelID, "Hello! I'm Gomez!")
	})
	router.Add("ping", func() {
		session.ChannelMessageSend(messageCreate.Message.ChannelID, "Pong!")
	})
	router.Add("info", func() {
		if len(messageCreate.Message.Mentions) == 0 {
			return
		}
		mentioned := messageCreate.Message.Mentions[0]
		session.ChannelMessageSend(messageCreate.Message.ChannelID, fmt.Sprint(mentioned))
	})
	router.Add("button", func() {
		buttonsJoined := router.RemovePrefixAndCommand(messageCreate.Message.Content)
		buttons := strings.Split(buttonsJoined, " ")
		session.ChannelMessageSendWithButtons(messageCreate.Message.ChannelID,"Here are your buttons:", buttons...)
		session.Client.On("messageButtonClicked", func() {
			session.ChannelMessageSend(messageCreate.Message.ChannelID, "Someone pressed the button!")
		})
	})
	router.Add("help", func() {
		var commands string
		for command,_ := range router.Routes {
			commands = commands + " " + command
		}
		session.ChannelMessageSend(messageCreate.Message.ChannelID, "Here are my current commands: "+commands)
	})
	router.Route(messageCreate.Message.Content)
}
