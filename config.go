package telegram

import (
	"encoding/json"
	"os"
	"time"
)

var (
	botLog   = NewLog(os.Stdout)
	botCache = make(map[string]*telegramBot)
)

type Config struct {
	Alias    string
	Token    string
	Webhook  string
	MsgStore Store
}

type telegramBot struct {
	*messageQueue
}

func RegisterBot(config *Config) error {
	bot := &telegramBot{}

	if config.Token == "" || config.MsgStore == nil {
		return NewError(InvalidConfig)
	}

	bot.client = newBotClient(config.Token, config.Webhook)

	if config.MsgStore != nil {
		bot.store = config.MsgStore
		bot.size = 10
		bot.start()
		//bot.queue = newMessageQueue(config.MsgStore)
	}

	botCache[config.Alias] = bot
	return nil
}

func newBot() *telegramBot {
	return newBotUsing(`default`)
}

func newBotUsing(alias string) *telegramBot {
	return botCache[alias]
}

func (b *telegramBot) ProcessMessage(message *Message) error {
	return b.client.processMessage(message)
}
func (b *telegramBot) PushMessage(message string) error {
	return b.store.RPush(message)
}

func ProcessMessage(message *Message) error {
	return newBot().client.processMessage(message)
}

func PushTextMessage(chatId int64, messageId int, message string) error {
	msg := &telegramMessage{
		ChatId:        chatId,
		MessageId:     messageId,
		Message:       message,
		Type:          MessageTypeText,
		RetryCount:    0,
		RetryInterval: RetryInterval,
		NextTime:      time.Now(),
	}
	bytes, _ := json.Marshal(msg)
	return newBot().PushMessage(string(bytes))
}

func PushPhotoMessage(chatId int64, messageId int, imgUrl, caption string) error {
	msg := &telegramMessage{
		ChatId:        chatId,
		MessageId:     messageId,
		Type:          MessageTypePhoto,
		ImgUrl:        imgUrl,
		Caption:       caption,
		RetryCount:    0,
		RetryInterval: RetryInterval,
		NextTime:      time.Now(),
	}
	bytes, _ := json.Marshal(msg)
	return newBot().PushMessage(string(bytes))
}
