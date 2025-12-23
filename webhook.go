package telegram

//func (b *TelegramBot) HandleWebhook(c *gin.Context) {
//
//	if c.Request.Method != "POST" {
//		http.Error(c.Writer, "[TelegramBot.HandleWebhook] Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//	var update Update
//	if err := c.ShouldBindJSON(&update); err != nil {
//		http.Error(c.Writer, "[TelegramBot.HandleWebhook] Failed to read request body", http.StatusBadRequest)
//		return
//	}
//
//	marshal, _ := json.Marshal(update)
//	botLog.Println("[TelegramBot.HandleWebhook] Received update %+s", string(marshal))
//
//	if update.Message == nil {
//		return
//	}
//
//	botClient := NewBot()
//	if botClient == nil {
//		botLog.Println("[TelegramBot.HandleWebhook] Failed to initialize bot client")
//		return
//	}
//
//	// Process the update
//	if err := botClient.ProcessMessage(update.Message); err != nil {
//		botLog.Printf("[TelegramBot.HandleWebhook] 处理消息异常：  err : %v \n", err)
//		return
//	}
//
//	// Respond to Telegram
//	c.String(http.StatusOK, "OK")
//}
