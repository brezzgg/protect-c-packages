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
	Pipes    []*Pipe
	TypeConv TypeConverter
}

func NewLogger(pipes []*Pipe, typeConv TypeConverter) *Logger {
	return &Logger{
		Pipes:    pipes,
		TypeConv: typeConv,
	}
}

func DefaultConsoleLogger() *Logger {
	return NewLogger(
		[]*Pipe{DefaultConsolePipe(false)},
		DefaultTypeConverter{},
	)
}

var pipeMutex sync.Mutex

func (l *Logger) Handle(m []any, level LogLevel, ch chan<- any) {
	pipeMutex.Lock()
	go func(m []any, level LogLevel, cal Caller, time time.Time) {
		defer close(ch)
		var msg Message

		// try to find formatted error in args
		for i, item := range m {
			if e, ok := item.(error); ok && strings.HasPrefix(e.Error(), FormatedErrorKey) {
				formStr := strings.TrimPrefix(e.Error(), FormatedErrorKey)
				err := json.Unmarshal([]byte(formStr), &msg)
				if err == nil {
					// mark the level as formatted error
					msg.Level = level.AppendLevels("(formatted)")

					// put all formatted context to key 'formatted'
					if msg.Context != nil {
						temp := msg.Context
						msg.Context = make(C)
						msg.Context["formatted"] = temp
					}

					// remove formatted error argument
					m = append(m[:i], m[i+1:]...)
				}
			}
		}

		// handle arguments
		err := l.CheckArgs(m, &msg.Text, &msg.Context)
		if err != nil {
			msg = Message{
				Caller:  cal,
				Level:   logLevelLoggerError,
				Text:    "Failed to parse message context",
				Context: C{"parserError": err.Error()},
			}
		} else {
			if msg.Caller.Equal(Caller{}) {
				msg.Caller = cal
			}
			if msg.Level.Equal(LogLevel{}) {
				msg.Level = level
			}
		}

		msg.Time = time

		var wg sync.WaitGroup
		wg.Add(len(l.Pipes))

		for _, pipe := range l.Pipes {
			go func(pipe *Pipe) {
				pipe.Handle(msg)
				wg.Done()
			}(pipe)
		}

		wg.Wait()
		pipeMutex.Unlock()

		ch <- 1
	}(m, level, GetCallerInfo(2), time.Now())
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
