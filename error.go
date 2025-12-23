package telegram

const (
	InvalidConfig         = 10400
	RequestFailure        = 10401 //请求失败
	TelegramApiError      = 10402
	TelegramBotError      = 10410
	SendMessageError      = 10411
	ParseResponseError    = 10412
	CommandNotFoundError  = 10413
	IllegalParameterError = 10414
	MessageTypeError      = 10415
)

var errorMessage = map[int]string{
	InvalidConfig:         "Invalid bot config",
	RequestFailure:        "request failure",
	TelegramApiError:      "telegram API error",
	TelegramBotError:      "telegram bot is nil",
	SendMessageError:      "Telegram send message error",
	ParseResponseError:    "failed to parse response",
	CommandNotFoundError:  "command not found",
	IllegalParameterError: "illegal parameter",
	MessageTypeError:      "Unsupported message type",
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err *Error) ErrorCode() int {
	return err.Code
}
func (err *Error) Error() string {
	return err.Msg
}

func NewError(code int, msg ...string) *Error {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = errorMessage[code]
	}

	return &Error{
		Code: code,
		Msg:  message,
	}
}
