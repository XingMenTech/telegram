package telegram

import (
	"io"
	"log"
)

type Log struct {
	*log.Logger
}

func NewLog(out io.Writer) *Log {
	d := new(Log)
	d.Logger = log.New(out, "[TG]", log.LstdFlags)
	return d
}
