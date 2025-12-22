package telegram

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var BotLog = NewLog(os.Stdout)

var botCache = make(map[string]*Bot)

var DefaultToken = ""

func RegisterBot(token string, isDefault bool, params ...Options) error {
	if _, ok := botCache[token]; ok {
		return nil
	}
	if isDefault {
		DefaultToken = token
	}
	ops := []Options{
		WithToken(token),
		WithParse(NewCommandParser("/")),
	}
	ops = append(ops, params...)
	bot, err := NewBotWidthOptions(ops...)
	if err != nil {
		return err
	}
	botCache[token] = bot
	return nil
}

func NewBot(token ...string) *Bot {
	if len(token) == 0 {
		return botCache[DefaultToken]
	}
	return botCache[token[0]]
}

// Bot represents a Telegram bot client
type Bot struct {
	token   string
	baseURL string
	parse   *CommandParser
	client  *http.Client
}

type Options func(*Bot) error

func WithToken(token string) Options {
	return func(b *Bot) error {
		b.token = token
		b.baseURL = fmt.Sprintf("https://api.telegram.org/bot%s", token)
		return nil
	}
}

func WithHook(webhook string) Options {
	return func(b *Bot) error {
		if err := b.SetWebhook(webhook); err != nil {
			log.Printf("telegram机器人初始化失败，errpr:%+v", err)
			return err
		}
		return nil
	}
}

func WithParse(parse *CommandParser) Options {
	return func(b *Bot) error {
		b.parse = parse
		return nil
	}
}

func NewBotWidthOptions(ops ...Options) (*Bot, error) {

	options := &Bot{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, op := range ops {
		if err := op(options); err != nil {
			BotLog.Printf("telegram机器人初始化失败，errpr:%v \n", err)
			return nil, err
		}
	}

	return options, nil
}
