package telegram

import (
	"os"
)

var (
	botLog   = NewLog(os.Stdout)
	botCache = make(map[string]*TelegramBot)
)

type Config struct {
	Alias string       `json:"alias" yaml:"alias"`
	Bot   *BotConfig   `json:"bot" yaml:"bot"`
	Queue *QueueConfig `json:"queue" yaml:"queue"`
}
type BotConfig struct {
	Token   string `json:"token" yaml:"token"`
	Webhook string `json:"webhook" yaml:"webhook"`
}
type QueueConfig struct {
	Type  string       `json:"type" yaml:"type"`
	Redis *RedisConfig `json:"redis" yaml:"redis"`
}
type RedisConfig struct {
	Host     string `json:"host" yaml:"host"`
	Password string `json:"password" yaml:"password"`
	DbNum    int    `json:"dbNum" yaml:"dbNum"`
}

type TelegramBot struct {
	client *BotClient
	queue  *MessageQueue
}

func RegisterBot(config *Config) error {
	bot := &TelegramBot{}

	botConfig := config.Bot
	if botConfig == nil {
		return NewError(InvalidConfig)
	}

	bot.client = NewBotClient(botConfig.Token, botConfig.Webhook)

	queueConfig := config.Queue
	if queueConfig != nil {
		var store Store
		if queueConfig.Type == "redis" {
			redisConfig := queueConfig.Redis
			store = NewRedisStore(redisConfig.Host, redisConfig.Password, redisConfig.DbNum)
		} else {
			store = NewListStore()
		}

		bot.queue = NewMessageQueue(store)
	}

	botCache[config.Alias] = bot
	return nil
}

func NewBot() *TelegramBot {
	return NewBotUsing(`default`)
}

func NewBotUsing(alias string) *TelegramBot {
	return botCache[alias]
}

func (b *TelegramBot) ProcessMessage(message *Message) error {
	return b.client.ProcessMessage(message)
}
