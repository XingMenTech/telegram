package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {

	if c.Request.Method != "POST" {
		http.Error(c.Writer, "[TelegramBot.HandleWebhook] Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var update Update
	if err := c.ShouldBindJSON(&update); err != nil {
		http.Error(c.Writer, "[TelegramBot.HandleWebhook] Failed to read request body", http.StatusBadRequest)
		return
	}

	marshal, _ := json.Marshal(update)
	BotLog.Println("[TelegramBot.HandleWebhook] Received update %+s", string(marshal))

	if update.Message == nil {
		return
	}

	botClient := NewBot()
	if botClient == nil {
		BotLog.Println("[TelegramBot.HandleWebhook] Failed to initialize bot client")
		return
	}

	// Process the update
	if err := botClient.ProcessMessage(update.Message); err != nil {
		BotLog.Printf("[TelegramBot.HandleWebhook] 处理消息异常：  err : %v \n", err)
		return
	}

	// Respond to Telegram
	c.String(http.StatusOK, "OK")
}

// SetWebhook sets the webhook URL for the bot
func (b *Bot) SetWebhook(url string) error {

	params := map[string]interface{}{
		"url": url,
	}
	respBody, err := b.doRequest("setWebhook", params)
	if err != nil {
		BotLog.Printf("[TelegramBot.SetWebhook] 设置webhook异常：  err : %v \n", err)
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		BotLog.Printf("[TelegramBot.SetWebhook] failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("[TelegramBot.SetWebhook] telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}

	return nil
}

// DeleteWebhook removes the webhook integration
func (b *Bot) DeleteWebhook() error {
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
		BotLog.Printf("[TelegramBot.DeleteWebhook] failed to parse response: %+v \n", err)
		return NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("[TelegramBot.DeleteWebhook] telegram API error: %v \n", result["description"])
		return NewError(TelegramApiError)
	}
	return nil
}

// GetWebhookInfo gets current webhook status
func (b *Bot) GetWebhookInfo() (map[string]interface{}, error) {

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
		BotLog.Printf("[TelegramBot.GetWebhookInfo] failed to parse response: %+v \n", err)
		return nil, NewError(ParseResponseError)
	}

	if !result["ok"].(bool) {
		BotLog.Printf("[TelegramBot.GetWebhookInfo] telegram API error: %v \n", result["description"])
		return nil, NewError(TelegramApiError)
	}
	return result["result"].(map[string]interface{}), nil
}
