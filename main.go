package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/khusainnov/vk-bot/bot"
	"go.uber.org/zap"
)

type config struct {
	Telegram struct {
		Token string `env:"TOKEN" envDefault:"6262174200:AAH6a9pg_JgtfBfs2EsMVAnbNBwYq5B4Fnw"`
	} `envPrefix:"VK_TG_"`
}

func main() {
	log, _ := zap.NewProduction()
	cfg := &config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatal("cannot parse config", zap.Error(err))
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("cannot init tgbotapi", zap.Error(err))
	}

	tgListener := &bot.TelegramListener{
		L:      log,
		BotAPI: botAPI,
	}

	log.Info("vk-test-bot started", zap.String("username", botAPI.Self.UserName))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err = tgListener.Do(ctx); err != nil {
		log.Info("telegram listener exit", zap.Error(err))
	}
}
