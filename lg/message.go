package lg

import "time"

type Message struct {
	Time    time.Time `json:"time"`
	Caller  Caller    `json:"caller"`
	Level   LogLevel  `json:"level"`
	Text    string    `json:"msg"`
	Context C         `json:"ctx,omitempty"`
}

func NewMessage(text string, level LogLevel, caller Caller, cont C) Message {
	return Message{
		Time:    time.Now(),
		Caller:  caller,
		Level:   level,
		Text:    text,
		Context: cont,
	}
}
