package main

import (
	"flag"
	"log"

	"inoutlog/telegram/client"
	"inoutlog/telegram/server"
	"inoutlog/timelogger"
)

func main() {
	token := flag.String("token", "", "telegram bot token see: https://core.telegram.org/api#bot-api")
	recordsFile := flag.String("out", "out.json", "where the output json is saved to")
	tariff := flag.Int("tariff", 50, "the hourly tariff")
	extra := flag.Int("extra", 30, "added to each record payment")
	flag.Parse()

	tclient, err := client.NewClient(*token)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := timelogger.NewLogger(*tariff, *extra, *recordsFile)
	if err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(tclient, logger)
	for {
		updates, err := s.TelegramClient.GetUpdates(60)
		if err != nil {
			log.Println("Something went wrong getting updates: " + err.Error())
		}
		err = s.HandleUpdates(updates)
		if err != nil {
			log.Println("Something went wrong handling updates: " + err.Error())
		}
	}
}
