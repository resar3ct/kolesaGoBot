package main

import (
	"flag"
	"kolesaGoBot/cmd/bot"
	"kolesaGoBot/cmd/message"

	"kolesaGoBot/internal/models"
	"log"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Env      string
	BotToken string
	Dsn      string
}

func main() {
	// go message.RunServer(":8001")

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := &Config{}
	_, err := toml.DecodeFile(*configPath, cfg)

	if err != nil {
		log.Fatalf("Ошибка декодирования файла конфигов %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	kolesaBot := bot.KolesaBot{
		Bot:   bot.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
	}

	messageBot := message.UpgradeBot{
		Bot:   kolesaBot.Bot,
		Users: &models.UserModel{Db: db},
	}
	go messageBot.RunServer(":8001")

	kolesaBot.Bot.Handle("/start", kolesaBot.StartHandler)

	kolesaBot.Bot.Start()
}
