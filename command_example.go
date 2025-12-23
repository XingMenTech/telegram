package telegram

import (
	"fmt"
	"strings"
)

// Example of how to use the command parser

func ExampleCommandParser() {
	// Create a new bot instance (you would use your actual bot token)
	// bot, err := NewBot("YOUR_BOT_TOKEN")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create a new command parser with "/" as the command prefix
	parser := NewCommandParser("/")

	// Register a simple "hello" command
	RegisterCommandFunc("hello", func(bot *BotClient, command *Command) error {
		response := fmt.Sprintf("Hello, %s!", command.Message.From.FirstName)
		return bot.SendMessage(command.Message.Chat.ID, response)
	})

	// Register a "help" command
	//RegisterCommandFunc("help", parser.DefaultHelpHandler())

	// Register a "echo" command with arguments
	RegisterCommandFunc("echo", func(bot *BotClient, command *Command) error {
		if len(command.Arguments) == 0 {
			return bot.SendMessage(command.Message.Chat.ID, "Usage: /echo <text>")
		}

		response := strings.Join(command.Arguments, " ")
		return bot.SendMessage(command.Message.Chat.ID, response)
	})

	// Add middleware to log commands
	//Use(func(bot *BotClient, command *Command) error {
	//	log.Printf("Processing command: %s, Arguments: %v", command.Name, command.Arguments)
	//	// Return true to continue processing
	//	return nil
	//})

	// Add middleware to restrict access to certain users
	//parser.Use(func(bot *BotClient, command *Command) bool {
	//	// Only allow user with ID 123456 to execute commands
	//	// Remove or modify this for your needs
	//	if command.Name == "admin" && command.Message.From.ID != 123456 {
	//		bot.SendMessage(command.Message.Chat.ID, "You don't have permission to use this command!")
	//		return false // Stop processing
	//	}
	//	return true // Continue processing
	//})

	// Example of processing an update
	// In a real application, this would be called from your webhook handler
	/*
		update := &Update{
			// ... update data from Telegram ...
		}
		err = parser.ProcessUpdate(bot, update)
		if err != nil {
			log.Printf("Error processing update: %v", err)
		}
	*/

	// Or process a message directly
	/*
		message := &Message{
			// ... message data ...
		}
		err = parser.ProcessMessage(bot, message)
		if err != nil {
			log.Printf("Error processing message: %v", err)
		}
	*/
	_ = parser
}
