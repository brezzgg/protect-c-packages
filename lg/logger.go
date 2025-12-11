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

var endTasks = &EndTasks{}

type (
	Logger struct {
		typeConv TypeConverter
		pipes    []*Pipe

		queue     chan Message
		wg        sync.WaitGroup
		stop      chan struct{}
		closed    atomic.Bool
		levelOpts []LogLevelOption
	}

	LoggerOption func(*Logger)
)

func NewLogger(opts ...LoggerOption) *Logger {
	l := &Logger{
		pipes:    make([]*Pipe, 0, 1),
		typeConv: NewDefaultTypeConverter(),
		queue:    make(chan Message, 8192),
		stop:     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(l)
	}

	if len(l.pipes) == 0 {
		l.pipes = append(l.pipes, NewPipe())
	}

	l.wg.Add(1)
	go l.worker()
	go l.watchSignals()

	l.closed.Store(false)

	return l
}

func WithPipe(pipe *Pipe) LoggerOption {
	return func(l *Logger) {
		l.pipes = append(l.pipes, pipe)
	}
}

func WithPipes(pipes ...*Pipe) LoggerOption {
	return func(l *Logger) {
		l.pipes = append(l.pipes, pipes...)
	}
}

func WithCustomTypeConverter(c TypeConverter) LoggerOption {
	return func(l *Logger) {
		l.typeConv = c
	}
}

func WithQueueSize(size int) LoggerOption {
	return func(l *Logger) {
		l.queue = make(chan Message, size)
	}
}

func WithConstantLevelOptions(opts ...LogLevelOption) LoggerOption {
	return func(l *Logger) {
		l.levelOpts = opts
	}
}

func (l *Logger) Close() {
	if !l.closed.Swap(true) {
		close(l.queue)
		close(l.stop)
		l.wg.Wait()
		for _, pipe := range l.pipes {
			pipe.Close()
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
				for _, p := range l.pipes {
					p.Handle(msg)
				}
			}
			break WorkerLoop

		case msg := <-l.queue:
			if msg.Text == "" {
				continue
			}
			for _, p := range l.pipes {
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

	// TODO: fix
	offset := 0
	caller := GetCallerInfo(2)
	if strings.HasPrefix(caller.File, "github.com/brezzgg/go-packages/lg") {
		offset = 1
	}

	msg := l.buildMessage(args, level, offset)
	l.queue <- msg
}

func (l *Logger) buildMessage(m []any, level LogLevel, callerOffset int) Message {
	level.opts = append(l.levelOpts, level.opts...)

	msg := Message{
		Caller: GetCallerInfo(3 + callerOffset),
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
		l.typeConv.ConvAndPush(item, func(k string, v any) {
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
