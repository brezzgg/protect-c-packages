package lg

import (
	"io"
	"os"
)

type (
	consoleWriter struct {
		stdout  *os.File
		discard bool
	}

	ConsoleWriterOption func(*consoleWriter)
)

func NewConsoleWriter(opts ...ConsoleWriterOption) *consoleWriter {
	w := &consoleWriter{
		stdout: os.Stdout,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func WithCustomStdout(f *os.File) ConsoleWriterOption {
	return func(w *consoleWriter) {
		w.stdout = f
	}
}

func WithWriterDiscard() ConsoleWriterOption {
	return func(w *consoleWriter) {
		w.discard = true
	}
}

func (c *consoleWriter) Write(str string) error {
	if c.discard {
		_, _ = io.Discard.Write([]byte(str))
		return nil
	}

	if _, err := c.stdout.WriteString(str + "\n"); err != nil {
		return err
	}

	return nil
}

func (c *consoleWriter) Flush() {
	_ = c.stdout.Sync()
}

func (c *consoleWriter) Close() {
	c.Flush()
}
