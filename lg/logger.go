package lg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Logger struct {
	Pipes     []*Pipe
	TypeConv  TypeConverter
	EndTasks  *EndTasks
	pipeMutex sync.RWMutex
}

func NewLogger(pipes []*Pipe, typeConv TypeConverter) *Logger {
	return &Logger{
		Pipes:    pipes,
		TypeConv: typeConv,
		EndTasks: &EndTasks{},
	}
}

func DefaultConsoleLogger() *Logger {
	return NewLogger(
		[]*Pipe{DefaultConsolePipe(false)},
		DefaultTypeConverter{},
	)
}

func (l *Logger) Handle(m []any, level LogLevel, ch chan<- any) {
	mCopy := make([]any, len(m))
	copy(mCopy, m)

	caller := GetCallerInfo(2)
	timestamp := time.Now()

	go func() {
		defer func() {
			if ch != nil {
				close(ch)
			}
		}()

		var msg Message

		filtered := make([]any, 0, len(mCopy))
		for _, item := range mCopy {
			if e, ok := item.(error); ok && strings.HasPrefix(e.Error(), FormatedErrorKey) {
				formStr := strings.TrimPrefix(e.Error(), FormatedErrorKey)
				if err := json.Unmarshal([]byte(formStr), &msg); err == nil {
					msg.Level = level.AppendLevels("(formatted)")

					if msg.Context != nil {
						temp := msg.Context
						msg.Context = C{"formatted": temp}
					}
					continue
				}
			}
			filtered = append(filtered, item)
		}

		if err := l.CheckArgs(filtered, &msg.Text, &msg.Context); err != nil {
			msg = Message{
				Caller:  caller,
				Level:   logLevelLoggerError,
				Text:    "Failed to parse message context",
				Context: C{"parserError": err.Error()},
				Time:    timestamp,
			}
		} else {
			if msg.Caller.Equal(Caller{}) {
				msg.Caller = caller
			}
			if msg.Level.Equal(LogLevel{}) {
				msg.Level = level
			}
			msg.Time = timestamp
		}

		l.pipeMutex.RLock()
		pipes := make([]*Pipe, len(l.Pipes))
		copy(pipes, l.Pipes)
		l.pipeMutex.RUnlock()

		if len(pipes) == 0 {
			return
		}

		for _, pipe := range pipes {
			go func(p *Pipe) {
				p.Handle(msg)
			}(pipe)
		}
	}()
}
func (l *Logger) CheckArgs(args []any, msgBody *string, msgCtx *C) error {
	appendBodyFunc := func(body string) {
		if len(*msgBody) > 0 {
			r := []rune((*msgBody))
			if unicode.IsUpper(r[0]) && unicode.IsLower(r[1]) {
				*msgBody = string(unicode.ToLower(r[0])) + (*msgBody)[1:]
			}
		}
		if len(*msgBody) == 0 {
			*msgBody = body
		} else {
			*msgBody = strings.TrimSuffix(strings.TrimSpace(body), ":") + ": " + strings.TrimSpace(*msgBody)
		}
	}

	if len(args) == 0 {
		return errors.New("no arguments")
	}
	if len(args) == 1 {
		if body, ok := args[0].(string); ok {
			appendBodyFunc(body)
			return nil
		}
		return errors.New("invalid argument")
	}

	if body, ok := args[0].(string); ok {
		appendBodyFunc(body)
	} else {
		return errors.New("invalid argument")
	}

	if *msgCtx == nil {
		*msgCtx = make(C)
	}
	for itemNum, item := range args {
		if itemNum == 0 {
			continue
		}
		l.TypeConv.ConvAndPush(item, func(k string, v any) {
			(*msgCtx)[l.findKey(k, *msgCtx)] = v
		})
	}

	return nil
}

func (*Logger) findKey(k string, m C) string {
	if _, ok := m[k]; !ok {
		return k
	}
	for i := 2; i < 100; i++ {
		keyStr := k + strconv.Itoa(i)
		if _, ok := m[keyStr]; !ok {
			return keyStr
		}
	}
	panic(fmt.Sprintf("too many duplicates of key '%s'", k))
}
