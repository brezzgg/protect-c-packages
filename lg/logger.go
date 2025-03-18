package lg

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
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
	go func(m []any, level LogLevel, cal Caller) {
		defer close(ch)
		var msg Message

		text, msgCtx, err := l.CheckArgs(m)
		if err != nil {
			msg = NewMessage("Failed to parse message context",
				logLevelLoggerError,
				cal,
				C{"parserError": err.Error()},
			)
		} else {
			msg = NewMessage(text,
				level,
				cal,
				msgCtx,
			)
		}

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
	}(m, level, GetCallerInfo(2))
}

func (l *Logger) CheckArgs(args []any) (string, C, error) {
	if len(args) == 0 {
		return "", C{}, errors.New("no arguments")
	}
	if len(args) == 1 {
		if msg, ok := args[0].(string); ok {
			return msg, C{}, nil
		}
		return "", C{}, errors.New("invalid argument")
	}

	var msg string
	if r, ok := args[0].(string); ok {
		msg = r
	} else {
		return "", C{}, errors.New("invalid argument")
	}

	msgCtx := make(C)
	for itemNum, item := range args {
		if itemNum == 0 {
			continue
		}
		l.TypeConv.ConvAndPush(item, func(k string, v any) {
			msgCtx[l.findKey(k, msgCtx)] = v
		})
	}

	return msg, msgCtx, nil
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
