package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
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
		fmt.Println("error: ", err)
		return
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	//イベントハンドラを追加
	discord.AddHandler(onMessageCreate)
	if err := discord.Open(); err != nil {
		fmt.Println(err)
		return
	}

	// 直近の関数（main）の最後に実行される
	defer discord.Close()

	fmt.Println("Listening...")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot自身のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		return
	}

	guildID := m.GuildID
	fmt.Println("guildID: ", guildID)

	// 特定のチャンネルのメッセージのみ処理
	// if m.ChannelID != "1295696391758286872" {
	// 	return
	// }

	// メンション対象のユーザーIDを設定
	// TODO: これをDBに格納する
	targetUserID := "1235248265276948505"

	for _, mention := range m.Mentions {
		fmt.Println("mentionID: ", mention.ID)
		fmt.Println("mentionName: ", mention.Username)
	}

	// メッセージ内容に特定のメンションが含まれるか確認
	if containsMention(m.Content, targetUserID) {
		// メッセージへの返信
		sendReply(s, m.ChannelID, "こんにちは！", m.Reference())
	}
}

func sendMessage(s *discordgo.Session, channelID string, msg string) {
	if _, err := s.ChannelMessageSend(channelID, msg); err != nil {
		fmt.Println("error: ", err)
	}
}

func sendReply(s *discordgo.Session, channelID string, msg string, reference *discordgo.MessageReference) {
	if _, err := s.ChannelMessageSendReply(channelID, msg, reference); err != nil {
		fmt.Println("error: ", err)
	}
}

// メッセージ内容に特定のユーザーへのメンションが含まれているか確認する関数
func containsMention(content string, userID string) bool {
	// メンションの形式を作成
	mention1 := "<@" + userID + ">"
	mention2 := "<@!" + userID + ">"

	return strings.Contains(content, mention1) || strings.Contains(content, mention2)
}
