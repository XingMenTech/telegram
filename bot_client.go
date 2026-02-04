package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// botClient represents a Telegram bot client
type botClient struct {
	token   string
	baseURL string
	parse   *commandParser
	client  *http.Client
}

type clientOptions func(*botClient) error

func withToken(token string) clientOptions {
	return func(b *botClient) error {
		b.token = token
		b.baseURL = fmt.Sprintf("https://api.telegram.org/bot%s", token)
		return nil
	}
}

func withHook(webhook string) clientOptions {
	return func(b *botClient) error {
		if err := b.setWebhook(webhook); err != nil {
			log.Printf("telegram机器人初始化失败，errpr:%+v", err)
			return err
		}
		return nil
	}
}

func withParse(parse *commandParser) clientOptions {
	return func(b *botClient) error {
		b.parse = parse
		return nil
	}
}

func newBotWidthOptions(ops ...clientOptions) (*botClient, error) {

	options := &botClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, op := range ops {
		if err := op(options); err != nil {
			botLog.Printf("telegram机器人初始化失败，errpr:%v \n", err)
			return nil, err
		}
	}

	return options, nil
}

func newBotClient(token, webhook string) *botClient {
	ops := []clientOptions{
		withToken(token),
		withParse(newCommandParser("/")),
	}
	if webhook != "" {
		ops = append(ops, withHook(webhook))
	}

	bot, err := newBotWidthOptions(ops...)
	if err != nil {
		return nil
	}
	return bot
}

// SendMessage sends a message to a chat
func (b *botClient) sendMessage(chatID int64, messageId int, text string) error {

	params := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	if messageId > 0 {
		params["reply_parameters"] = map[string]interface{}{
			"message_id": messageId,
		}
	}
	respBody, err := b.doRequest("sendMessage", params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		botLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}

	return nil
}

// ReplyMessage replies to a message
func (b *botClient) replyMessage(chatId int64, messageId int, text string) error {

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
		botLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// SendPhoto sends a photo to a chat
func (b *botClient) sendPhoto(chatID int64, photoURL, caption string) error {

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
		botLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// ForwardMessage forwards a message from one chat to another
func (b *botClient) forwardMessage(chatID, fromChatID int64, messageID int) error {

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
		botLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// CopyMessage copies a message from one chat to another
func (b *botClient) copyMessage(chatID, fromChatID int64, messageID int) error {

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
		botLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// GetFile gets information about a file by its file_id
func (b *botClient) getFile(fileID string) (*File, error) {

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
		botLog.Printf("failed to parse response: %+v \n", err)
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Ok {
		botLog.Printf("telegram API error: %v \n", result.Result)
		return nil, fmt.Errorf("telegram API error: %v", result.Result)
	}

	return &result.Result, nil
}

// GetFileURL returns the full URL to download a file
func (b *botClient) getFileURL(file *File) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.token, file.FilePath)
}

// SendMediaGroup sends a group of photos as an album
func (b *botClient) sendMediaGroup(chatID int64, media []InputMedia) error {

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
		botLog.Printf("failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// GetUpdates retrieves updates from the bot
func (b *botClient) getUpdates(offset int64, limit int) ([]Update, error) {

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
		botLog.Printf("failed to parse response: %+v \n", err)
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !result.Ok {
		botLog.Printf("telegram API error: %v \n", result.Result)
		return nil, fmt.Errorf("telegram API error")
	}

	return result.Result, nil
}

func (b *botClient) processUpdate(update *Update) error {
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
		return NewError(CommandNotFoundError)
	}
	return command.Handler()
}

// ProcessMessage 处理消息并执行相应的命令处理程序
func (b *botClient) processMessage(message *Message) error {
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
		return NewError(CommandNotFoundError)
	}

	return command.Handler()
}

// SetWebhook sets the webhook URL for the bot
func (b *botClient) setWebhook(url string) error {

	params := map[string]interface{}{
		"url": url,
	}
	respBody, err := b.doRequest("setWebhook", params)
	if err != nil {
		botLog.Printf("[TelegramBot.SetWebhook] 设置webhook异常：  err : %v \n", err)
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		botLog.Printf("[TelegramBot.SetWebhook] failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("[TelegramBot.SetWebhook] telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}

	return nil
}

// DeleteWebhook removes the webhook integration
func (b *botClient) deleteWebhook() error {
	reqUrl := fmt.Sprintf("%s/deleteWebhook", b.baseURL)

	resp, err := b.client.Get(reqUrl)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		botLog.Printf("[TelegramBot.DeleteWebhook] failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("[TelegramBot.DeleteWebhook] telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// GetWebhookInfo gets current webhook status
func (b *botClient) getWebhookInfo() (map[string]interface{}, error) {

	url := fmt.Sprintf("%s/getWebhookInfo", b.baseURL)

	resp, err := b.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook info: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		botLog.Printf("[TelegramBot.GetWebhookInfo] failed to parse response: %+v \n", err)
		return nil, NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		botLog.Printf("[TelegramBot.GetWebhookInfo] telegram API error: %v \n", result["description"])
		return nil, NewError(TelegramApiError)
	}
	return result["result"].(map[string]interface{}), nil
}

func (b *botClient) doRequest(api string, params map[string]interface{}) (body []byte, err error) {

	reqUrl := fmt.Sprintf("%s/%s", b.baseURL, api)
	botLog.Printf("[TelegramBot.Request] 请求地址：%s \n", reqUrl)

	paramBytes, _ := json.Marshal(params)
	botLog.Printf("[TelegramBot.Request] 请求参数：%s", string(paramBytes))

	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(paramBytes))
	if err != nil {
		botLog.Printf("[TelegramBot.Request] 创建请求异常：  err : %v \n", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		botLog.Printf("[TelegramBot.Request] 发送请求异常：  err : %v \n", err)
		return nil, fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		botLog.Printf("[TelegramBot.Request] request failed with status %d \n", resp.StatusCode)
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		botLog.Printf("[TelegramBot.Request] 读取响应异常：  err : %v \n", err)
	}
	botLog.Printf("[TelegramBot.Request] 响应参数：%s", string(body))
	return
}
