package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TelegramListener struct {
	L      *zap.Logger
	BotAPI *tgbotapi.BotAPI
}

func (t *TelegramListener) GreetingMessage(msg *tgbotapi.Message) error {
	rspMsg := "Привет! Я vk_bot.\n\nХочу рассказать тебе о стажировках которые есть в компании VK.\nНажав на `меню` – тебе откроется список возможных команд\n\n"

	menuKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Меню", menuCmd),
		),
	)

	rsp := tgbotapi.NewMessage(msg.Chat.ID, rspMsg)
	rsp.ReplyMarkup = menuKeyboard

	if _, err := t.BotAPI.Send(rsp); err != nil {
		return fmt.Errorf("error due send message, %w", err)
	}

	return nil
}

func (t *TelegramListener) Do(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updChan := t.BotAPI.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			t.L.Error("called context done, listener will stop", zap.Error(ctx.Err()))
			t.BotAPI.StopReceivingUpdates()
			return ctx.Err()
		case upd, ok := <-updChan:
			if !ok {
				t.L.Error("closed update channel")
				return errClosedCh
			}

			if upd.Message != nil {
				if err := t.handleMessage(upd.Message); err != nil {
					t.L.Error("cannot handle message", zap.Error(err))
				}
			}

			if upd.CallbackQuery != nil {
				if err := t.handleCallback(upd.CallbackQuery); err != nil {
					t.L.Error("cannot handle callback", zap.Error(err))
				}
			}

		}
	}
}

func (t *TelegramListener) handleMessage(msg *tgbotapi.Message) error {
	if msg.Document != nil {
		return errUnsupported
	}
	if msg.Poll != nil {
		return errUnsupported
	}

	cmd := parseCommand(msg.Text)

	switch cmd {
	case startCmd:
		if err := t.GreetingMessage(msg); err != nil {
			return err
		}
	default:
		rspMsg := tgbotapi.NewMessage(msg.Chat.ID, "данной команды не существует")
		rspMsg.ReplyToMessageID = msg.MessageID
		if _, err := t.BotAPI.Send(rspMsg); err != nil {
			t.L.Error(errUnsupported.Error(), zap.Error(err))
			return fmt.Errorf("%v, %w", errUnsupported, err)
		}
	}

	return nil
}

func (t *TelegramListener) handleCallback(cb *tgbotapi.CallbackQuery) error {
	if cb.Message.Poll != nil {
		return errUnsupported
	}
	if cb.Message.Document != nil {
		return errUnsupported
	}

	cmd := parseCommand(cb.Data)

	switch cmd {
	case menuCmd:
		if err := t.sendMenu(cb); err != nil {
			return err
		}
	case whatForImCmd:
		rspMsg := tgbotapi.NewMessage(cb.Message.Chat.ID, whatForMsg)
		if _, err := t.BotAPI.Send(rspMsg); err != nil {
			return fmt.Errorf("cannot send what for i'm message, %w", err)
		}
	case internshipCmd:
		if err := t.internship(cb); err != nil {
			return err
		}
	case routesCmd:
		if err := t.routes(cb); err != nil {
			return err
		}
	default:
		rspMsg := tgbotapi.NewMessage(cb.Message.Chat.ID, "данной команды не существует")
		rspMsg.ReplyToMessageID = cb.Message.MessageID
		if _, err := t.BotAPI.Send(rspMsg); err != nil {
			t.L.Error(errUnsupported.Error(), zap.Error(err))
			return fmt.Errorf("%v, %w", errUnsupported, err)
		}
	}

	return nil
}

func (t *TelegramListener) sendMenu(cb *tgbotapi.CallbackQuery) error {
	msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "Мои команды:\n\n")

	menuKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Для чего я нужен", whatForImCmd),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расскажи о стажировке", internshipCmd),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Какие есть направления", routesCmd),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Подробнее о стажировке", "https://internship.vk.company/internship"),
		),
	)

	msg.ReplyMarkup = menuKeyboard

	if _, err := t.BotAPI.Send(msg); err != nil {
		return fmt.Errorf("cannot send menu message, %w", err)
	}

	return nil
}

func (t *TelegramListener) internship(cb *tgbotapi.CallbackQuery) error {
	rspMsg := tgbotapi.NewMessage(cb.Message.Chat.ID, internshipMsg)
	data, err := readFile("internship.png")
	if err != nil {
		return fmt.Errorf("cannot read image, %w", err)
	}

	imgData := tgbotapi.FileBytes{
		Name:  "internship",
		Bytes: data,
	}

	img := tgbotapi.NewPhoto(cb.Message.Chat.ID, imgData)
	if _, err = t.BotAPI.Send(img); err != nil {
		return fmt.Errorf("cannnot send internship image, %w", err)
	}

	if _, err = t.BotAPI.Send(rspMsg); err != nil {
		return fmt.Errorf("cannot send internship message, %w", err)
	}
	return nil
}

func (t *TelegramListener) routes(cb *tgbotapi.CallbackQuery) error {
	msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "Направления стажировки")

	routesKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Разработка", "https://internship.vk.company/internship?direction=5"),
			tgbotapi.NewInlineKeyboardButtonURL("Дизайн", "https://internship.vk.company/internship?direction=10"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Информационная безопасноть", "https://internship.vk.company/internship?direction=3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Продажи в IT", "https://internship.vk.company/internship?direction=17"),
			tgbotapi.NewInlineKeyboardButtonURL("Маркетинг", "https://internship.vk.company/internship?direction=15"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Аналитика", "https://internship.vk.company/internship?direction=6"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Управление проектами", "https://internship.vk.company/internship?direction=9"),
		),
	)

	msg.ReplyMarkup = routesKeyboard

	if _, err := t.BotAPI.Send(msg); err != nil {
		return fmt.Errorf("cannot sent routes message, %w", err)
	}

	return nil
}
