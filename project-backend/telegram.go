package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendToTelegram(filePath string) error {
	bot, err := tgbotapi.NewBotAPI("7952271289:AAG3CaGezk_Rwr4o9LdIOj1HkKU3nLVsxvU")
	if err != nil {
		log.Println("Telegram bot error:", err)
		return err
	}

	chatID := int64(1819220129) // ← ТВОЙ chat ID

	_, _ = bot.Send(tgbotapi.NewMessage(chatID, "📎 Yeni müqavilələr ZIP faylında:"))

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	zipFile := tgbotapi.FileBytes{
		Name:  "contracts.zip",
		Bytes: fileBytes,
	}

	_, err = bot.Send(tgbotapi.NewDocument(chatID, zipFile))
	if err != nil {
		log.Println("Telegram send error:", err)
	}
	return err
}
