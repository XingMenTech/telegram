package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SendMessage sends a message to a chat
func (b *Bot) SendMessage(chatID int64, text string) error {

	params := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	respBody, err := b.doRequest("sendMessage", params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}

	return nil
}

// ReplyMessage replies to a message
func (b *Bot) ReplyMessage(chatId int64, messageId int, text string) error {

	params := map[string]interface{}{
		"chat_id": chatId,
		"text":    text,
		"reply_parameters": map[string]interface{}{
			"message_id": messageId,
		},
	}
	respBody, err := b.doRequest("sendMessage", params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// SendPhoto sends a photo to a chat
func (b *Bot) SendPhoto(chatID int64, photoURL, caption string) error {

	params := map[string]interface{}{
		"chat_id": chatID,
		"photo":   photoURL,
	}

	if caption != "" {
		params["caption"] = caption
	}
	respBody, err := b.doRequest("sendPhoto", params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// ForwardMessage forwards a message from one chat to another
func (b *Bot) ForwardMessage(chatID, fromChatID int64, messageID int) error {

	params := map[string]interface{}{
		"chat_id":      chatID,
		"from_chat_id": fromChatID,
		"message_id":   messageID,
	}
	respBody, err := b.doRequest("forwardMessage", params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// CopyMessage copies a message from one chat to another
func (b *Bot) CopyMessage(chatID, fromChatID int64, messageID int) error {

	params := map[string]interface{}{
		"chat_id":      chatID,
		"from_chat_id": fromChatID,
		"message_id":   messageID,
	}
	respBody, err := b.doRequest("copyMessage", params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// GetFile gets information about a file by its file_id
func (b *Bot) GetFile(fileID string) (*File, error) {

	params := map[string]interface{}{
		"file_id": fileID,
	}
	respBody, err := b.doRequest("getFile", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Ok     bool `json:"ok"`
		Result File `json:"result"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Ok {
		BotLog.Printf("telegram API error: %v \n", result.Result)
		return nil, fmt.Errorf("telegram API error: %v", result.Result)
	}

	return &result.Result, nil
}

// GetFileURL returns the full URL to download a file
func (b *Bot) GetFileURL(file *File) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.token, file.FilePath)
}

// SendMediaGroup sends a group of photos as an album
func (b *Bot) SendMediaGroup(chatID int64, media []InputMedia) error {

	params := map[string]interface{}{
		"chat_id": chatID,
		"media":   media,
	}
	respBody, err := b.doRequest("sendMediaGroup", params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// GetUpdates retrieves updates from the bot
func (b *Bot) GetUpdates(offset int64, limit int) ([]Update, error) {

	params := map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	}
	respBody, err := b.doRequest("getUpdates", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("failed to parse response: %+v \n", err)
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Ok {
		BotLog.Printf("telegram API error: %v \n", result.Result)
		return nil, fmt.Errorf("telegram API error")
	}

	return result.Result, nil
}

func (b *Bot) ProcessUpdate(update *Update) error {
	// Only process message updates
	if update.Message == nil {
		return nil
	}

	// Only process text messages
	if update.Message.Text == "" {
		return nil
	}

	// Parse command from message
	command := b.parse.ParseCommand(update.Message.Text, update.Message)
	if command == nil {
		return nil
	}
	return command.Handler(b)
}

// ProcessMessage 处理消息并执行相应的命令处理程序
func (b *Bot) ProcessMessage(message *Message) error {
	// Only process text messages
	if message.Text == "" && len(message.Photo) == 0 {
		return nil
	}

	commandText := message.Text
	if commandText == "" {
		commandText = message.Caption
	}

	// Parse command from message
	command := b.parse.ParseCommand(commandText, message)
	if command == nil {
		return nil
	}

	return command.Handler(b)
}

func (b *Bot) doRequest(api string, params map[string]interface{}) (body []byte, err error) {

	reqUrl := fmt.Sprintf("%s/%s", b.baseURL, api)
	BotLog.Printf("[TelegramBot.Request] 请求地址：%s \n", reqUrl)

	paramStr, _ := json.Marshal(params)
	BotLog.Printf("[TelegramBot.Request] 请求参数：%s", paramStr)

	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(body))
	if err != nil {
		BotLog.Printf("[TelegramBot.Request] 创建请求异常：  err : %v \n", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		BotLog.Printf("[TelegramBot.Request] 发送请求异常：  err : %v \n", err)
		return nil, fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		BotLog.Printf("[TelegramBot.Request] 读取响应异常：  err : %v \n", err)
	}
	BotLog.Printf("[TelegramBot.Request] 响应参数：%s", string(body))
	return
}
