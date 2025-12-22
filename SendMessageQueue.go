package telegram

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/ants/v2"
)

const (
	telegramSendMessageKey = "telegram:send_message"
	RetryMaxTimesCount     = 3 //重试最大次数
	RetryTimeLimitNum      = 2
	RetryInterval          = 10 //重试间隔时间
)

var (
	redisClient *redis.Client
	ctx         context.Context
)

type telegramMessage struct {
	ChatId        int64     `json:"chatId"`
	Message       string    `json:"message"`
	RetryCount    int       `json:"callbackCount"` //重试计数
	CallbackRaw   string    `json:"-"`             //回调原串
	RetryInterval int       `json:"retryInterval"` //重试间隔
	NextTime      time.Time `json:"nextTime"`      //下一次重试的时间
}

func AutoSendMessage(max int, client *redis.Client) {
	redisClient = client

	// 创建一个容量为max的goroutine池
	p, _ := ants.NewPoolWithFunc(max, func(i interface{}) {
		cb := i.(*telegramMessage)
		// 失败则放入重试队列
		if err := sendMessage(cb); err == nil {
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
		_ = putRetryCache(cb, cb.CallbackRaw, true)
	}, ants.WithPreAlloc(true), ants.WithNonblocking(false))
	defer p.Release()

	for {
		cmd := redisClient.BLPop(ctx, 0, telegramSendMessageKey) // redisClient.BLPop(ctx, telegramSendMessageKey, 0)
		if cmd.Err() != nil {
			time.Sleep(time.Second * 3)
			continue
		}

		jsonStr := cmd.Val()[1]
		var cb telegramMessage
		if err := json.Unmarshal([]byte(jsonStr), &cb); err != nil {
			BotLog.Printf("[telegram_send_message] json Unmarshal, origin : %s ,err :%v \n", jsonStr, err)
			continue
		}

		if !time.Now().After(cb.NextTime) {
			_ = putRetryCache(&cb, jsonStr, false)
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
			continue
		}

		cb.CallbackRaw, cb.RetryInterval = jsonStr, RetryInterval
		_ = p.Invoke(&cb)
	}
}

func putRetryCache(cb *telegramMessage, str string, isTimeNow bool) (err error) {
	if !isTimeNow {
		return redisClient.RPush(ctx, telegramSendMessageKey, str).Err()
	}
	cbByte, err := json.Marshal(cb)
	if err != nil {
		return redisClient.RPush(ctx, telegramSendMessageKey, str).Err()
	}

	return redisClient.RPush(ctx, telegramSendMessageKey, string(cbByte)).Err()
}

func sendMessage(cb *telegramMessage) error {

	bot := NewBot()
	if bot == nil {
		BotLog.Println("[telegram_send_message] telegram bot is nil")
		return NewError(TelegramBotError)
	}

	return bot.SendMessage(cb.ChatId, cb.Message)
}
