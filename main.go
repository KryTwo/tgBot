package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"main/repository"
	"main/service"
	"time"
)

var (
	// глобальная переменная в которой храним токен
	telegramBotToken string
)

/*func init() {
	// принимаем на входе флаг -telegrambottoken
	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Parse()

	// без него не запускаемся
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
}*/

func main() {
	db := repository.NewDB()
	defer db.Close()
	/*out, err := repository.GetStatsMe(db, "KryTwo")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", out)
	fmt.Printf("%T\n", out)*/

	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	var reply string
	var lastUse int

	for update := range updates {
		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "make_me_pidor":
			reply = service.NewUser(db, update.Message.From.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msg)
		case "who_is_pidor":
			reply = service.Stats(db)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msg)
		case "wheel_of_fortune":
			if lastUse == time.Now().YearDay() {
				reply = "Сегодня пидор уже выбран, приходите завтра."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
				continue
			}
			reply = service.IncRandUser(db, update.Message.From.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msg)
			lastUse = time.Now().YearDay()
		case "i_am_not_a_pidor":
			reply = service.StatsMe(db, update.Message.From.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msg)
		}

		time.Sleep(1 * time.Second)
	}

}
