package lg

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unicode"
)

type Logger struct {
	Pipes    []*Pipe
	TypeConv TypeConverter
	EndTasks *EndTasks

	queue  chan Message
	wg     sync.WaitGroup
	stop   chan struct{}
	closed atomic.Bool
}

func NewLogger(pipes []*Pipe, typeConv TypeConverter) *Logger {
	l := &Logger{
		Pipes:    pipes,
		TypeConv: typeConv,
		EndTasks: &EndTasks{},
		queue:    make(chan Message, 8192),
		stop:     make(chan struct{}),
	}

	l.wg.Add(1)
	go l.worker()

	go l.watchSignals()

	return l
}

func DefaultConsoleLogger() *Logger {
	return NewLogger(
		[]*Pipe{DefaultConsolePipe(false)},
		DefaultTypeConverter{},
	)
}

func (l *Logger) Close() {
	if !l.closed.Swap(true) {
		close(l.queue)
		close(l.stop)
		l.wg.Wait()
		for _, pipe := range l.Pipes {
			pipe.Wri.Flush()
		}
	}
}

func (l *Logger) worker() {
WorkerLoop:
	for {
		select {
		case <-l.stop:
			for msg := range l.queue {
				if msg.Text == "" {
					continue
				}
				for _, p := range l.Pipes {
					p.Handle(msg)
				}
			}
			break WorkerLoop

		case msg := <-l.queue:
			if msg.Text == "" {
				continue
			}
			for _, p := range l.Pipes {
				p.Handle(msg)
			}
		}
	}

	l.wg.Done()
}

func (l *Logger) watchSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh
	l.Close()
}

func (l *Logger) Handle(args []any, level LogLevel) {
	if l.closed.Load() {
		return
	}
	msg := l.buildMessage(args, level)
	l.queue <- msg
}

func (l *Logger) buildMessage(m []any, level LogLevel) Message {
	msg := Message{
		Caller: GetCallerInfo(3),
		Level:  level,
		Time:   time.Now(),
	}

	if err := l.CheckArgs(m, &msg.Text, &msg.Context); err != nil {
		msg = Message{
			Caller:  msg.Caller,
			Level:   logLevelLoggerError,
			Text:    "Failed to parse message context",
			Context: C{"parserError": err.Error()},
			Time:    msg.Time,
		}
	}

	return msg
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
