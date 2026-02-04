package telegram

import (
	"strings"
)

// Command represents a parsed command from a Telegram message
type Command struct {
	Name      string
	Arguments []string
	RawText   string
	Message   *Message
}

func (c *Command) Handler() error {
	// Run middleware
	for _, m := range middleware {
		if err := m(c); err != nil {
			return err
		}
	}

	// Find and execute command handler
	if handler, exists := commands[c.Name]; exists {
		return handler(c)
	}

	return nil
}

// CommandHandlerFunc is a function type that implements CommandHandler
type CommandHandlerFunc func(command *Command) error

var (
	commands   = make(map[string]CommandHandlerFunc)
	middleware = make([]CommandHandlerFunc, 0)
)

// RegisterCommandFunc 为特定命令注册处理程序函数
func RegisterCommandFunc(name string, handler CommandHandlerFunc) {
	commands[name] = handler
}

// Use 将中间件添加到命令解析器中
func Use(m ...CommandHandlerFunc) {
	middleware = append(middleware, m...)
}

// CommandParser 负责解析Telegram消息中的命令
type commandParser struct {
	prefix string
}

// NewCommandParser 使用指定的命令前缀创建新的命令解析器
func newCommandParser(prefix string) *commandParser {
	return &commandParser{
		prefix: prefix,
	}
}

// ParseCommand 从消息文本中解析命令
func (cp *commandParser) ParseCommand(text string, message *Message) *Command {
	if !strings.HasPrefix(text, cp.prefix) {
		return nil
	}

	// Remove prefix
	cmdText := strings.TrimPrefix(text, cp.prefix)

	// Split by spaces
	parts := strings.Fields(cmdText)

	if len(parts) == 0 {
		return nil
	}

	command := &Command{
		Name:    strings.ToLower(parts[0]),
		RawText: text,
		Message: message,
	}

	if len(parts) > 1 {
		command.Arguments = parts[1:]
	}

	return command
}
