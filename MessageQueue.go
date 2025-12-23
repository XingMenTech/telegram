package telegram

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	RetryMaxTimesCount = 3 //重试最大次数
	RetryTimeLimitNum  = 2
	RetryInterval      = 10 //重试间隔时间
)

var (
	queue *MessageQueue
	qOnce sync.Once
)

type telegramMessage struct {
	ChatId        int64     `json:"chatId"`
	Message       string    `json:"message"`
	Type          string    `json:"type"`
	ImgUrl        string    `json:"imgUrl"`
	Caption       string    `json:"caption"`
	RetryCount    int       `json:"callbackCount"` //重试计数
	CallbackRaw   string    `json:"-"`             //回调原串
	RetryInterval int       `json:"retryInterval"` //重试间隔
	NextTime      time.Time `json:"nextTime"`      //下一次重试的时间
}

type MessageQueue struct {
	store Store
	size  int
}

func NewMessageQueue(store Store) *MessageQueue {
	qOnce.Do(func() {
		queue = &MessageQueue{
			store: store,
			size:  10,
		}
		go queue.Start()
	})
	return queue
}

func (mq *MessageQueue) Start() {

	botLog.Println("[telegram_send_message] telegram message queue start。。。。。。。。。")
	// 创建一个容量为max的goroutine池
	p, _ := ants.NewPoolWithFunc(mq.size, func(i interface{}) {
		cb := i.(*telegramMessage)
		// 失败则放入重试队列
		if err := mq.sendMessage(cb); err == nil {
			return
		}

		if cb.RetryCount >= RetryMaxTimesCount {
			return
		}
		if cb.RetryCount >= RetryTimeLimitNum {
			cb.NextTime = time.Now().Add(2 * time.Minute)
		} else {
			cb.NextTime = time.Now().Add(time.Duration((cb.RetryCount+1)*(cb.RetryCount+1)*cb.RetryInterval) * time.Second)
		}
		cb.RetryCount++
		_ = mq.putRetryCache(cb, cb.CallbackRaw, true)
	}, ants.WithPreAlloc(true), ants.WithNonblocking(false))
	defer p.Release()

	for {

		jsonStr, cmdErr := mq.store.BLPop()
		if cmdErr != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		if jsonStr == "" {
			time.Sleep(time.Second * 3)
			continue
		}

		var cb telegramMessage
		if err := json.Unmarshal([]byte(jsonStr), &cb); err != nil {
			botLog.Printf("[telegram_send_message] json Unmarshal, origin : %s ,err :%v \n", jsonStr, err)
			continue
		}

		if !time.Now().After(cb.NextTime) {
			_ = mq.putRetryCache(&cb, jsonStr, false)
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
			continue
		}

		cb.CallbackRaw, cb.RetryInterval = jsonStr, RetryInterval
		_ = p.Invoke(&cb)
	}
}

func (mq *MessageQueue) sendMessage(cb *telegramMessage) error {
	bot := NewBot()
	if bot == nil {
		botLog.Println("[telegram_send_message] telegram bot is nil")
		return NewError(TelegramBotError)
	}

	switch cb.Type {
	case MessageTypeText:
		return bot.client.SendMessage(cb.ChatId, cb.Message)
	case MessageTypePhoto:
		return bot.client.SendPhoto(cb.ChatId, cb.ImgUrl, cb.Caption)
	default:
		return NewError(MessageTypeError)
	}
}

func (mq *MessageQueue) putRetryCache(cb *telegramMessage, raw string, isTimeNow bool) error {
	if !isTimeNow {
		return mq.store.RPush(raw)
	}
	cbByte, err := json.Marshal(cb)
	if err != nil {
		return mq.store.RPush(raw)
	}
	return mq.store.RPush(string(cbByte))
}

func PushTextMessage(chatId int64, message string) error {
	msg := &telegramMessage{
		ChatId:        chatId,
		Message:       message,
		Type:          MessageTypeText,
		RetryCount:    0,
		RetryInterval: RetryInterval,
		NextTime:      time.Now(),
	}
	bytes, _ := json.Marshal(msg)
	return NewBot().queue.store.RPush(string(bytes))
}

func PushPhotoMessage(chatId int64, imgUrl, caption string) error {
	msg := &telegramMessage{
		ChatId:        chatId,
		Type:          MessageTypePhoto,
		ImgUrl:        imgUrl,
		Caption:       caption,
		RetryCount:    0,
		RetryInterval: RetryInterval,
		NextTime:      time.Now(),
	}
	bytes, _ := json.Marshal(msg)
	return NewBot().queue.store.RPush(string(bytes))
}

//func AutoSendMessage(max int, store Store) {
//	messageStore = store
//
//	// 创建一个容量为max的goroutine池
//	p, _ := ants.NewPoolWithFunc(max, func(i interface{}) {
//		cb := i.(*telegramMessage)
//		// 失败则放入重试队列
//		if err := sendMessage(cb); err == nil {
//			return
//		}
//
//		if cb.RetryCount >= RetryMaxTimesCount {
//			return
//		}
//		if cb.RetryCount >= RetryTimeLimitNum {
//			cb.NextTime = time.Now().Add(2 * time.Minute)
//		} else {
//			cb.NextTime = time.Now().Add(time.Duration((cb.RetryCount+1)*(cb.RetryCount+1)*cb.RetryInterval) * time.Second)
//		}
//		cb.RetryCount++
//		_ = putRetryCache(cb, cb.CallbackRaw, true)
//	}, ants.WithPreAlloc(true), ants.WithNonblocking(false))
//	defer p.Release()
//
//	for {
//
//		jsonStr, cmdErr := messageStore.BLPop()
//		if cmdErr != nil {
//			time.Sleep(time.Second * 3)
//			continue
//		}
//		if jsonStr == "" {
//			time.Sleep(time.Second * 3)
//			continue
//		}
//
//		var cb telegramMessage
//		if err := json.Unmarshal([]byte(jsonStr), &cb); err != nil {
//			botLog.Printf("[telegram_send_message] json Unmarshal, origin : %s ,err :%v \n", jsonStr, err)
//			continue
//		}
//
//		if !time.Now().After(cb.NextTime) {
//			_ = putRetryCache(&cb, jsonStr, false)
//			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
//			continue
//		}
//
//		cb.CallbackRaw, cb.RetryInterval = jsonStr, RetryInterval
//		_ = p.Invoke(&cb)
//	}
//}

//func putRetryCache(cb *telegramMessage, str string, isTimeNow bool) (err error) {
//	// 根据使用的存储类型选择不同的方法
//
//	if !isTimeNow {
//		return messageStore.RPush(str)
//	}
//	cbByte, err := json.Marshal(cb)
//	if err != nil {
//		return messageStore.RPush(str)
//	}
//	return messageStore.RPush(string(cbByte))
//}
//
//func sendMessage(cb *telegramMessage) error {
//	bot := NewBot()
//	if bot == nil {
//		botLog.Println("[telegram_send_message] telegram bot is nil")
//		return NewError(TelegramBotError)
//	}
//
//	switch cb.Type {
//	case MessageTypeText:
//		return bot.SendMessage(cb.ChatId, cb.Message)
//	case MessageTypePhoto:
//		return bot.SendPhoto(cb.ChatId, cb.ImgUrl, cb.Caption)
//	default:
//		return NewError(MessageTypeError)
//	}
//}
