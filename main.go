package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	token = "sample-token"
)

func main() {
	discordToken := "Bot " + token

	discord, err := discordgo.New(discordToken)
	if err != nil {
		fmt.Println("discordgo.New error: ", err)
		return
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	//イベントハンドラを追加
	discord.AddHandler(onMessageCreate)
	if err := discord.Open(); err != nil {
		fmt.Println("discord.Open error: ", err)
		return
	}

	// 直近の関数（main）の最後に実行される
	defer discord.Close()

	fmt.Println("Listening...")

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stopBot
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot自身のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		return
	}

	if _, err := s.ChannelMessageSendReply(m.ChannelID, "こんにちは！", m.Reference()); err != nil {
		fmt.Println("discordgo.Session.ChannelMessageSendReply error: ", err)
	}
}
