package main

/*
TODO:

. Исправление ошибок
. Проверка на существование папок
. Детальные логи
. Шифрование данных, полученных от пользователя


*/

import (
	"flag"
	"log"
	tgClient "main/clients/telegram"
	event_consumer "main/consumer/event-consumer"
	"main/events/telegram"
	"main/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

// 5297059869:AAEQ7sgCb1vrua_QPm8Rtjs6cdv2DddJ13A
func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
