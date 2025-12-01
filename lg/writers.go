package lg

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type ConsoleWriter struct{}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (c ConsoleWriter) Write(str string) error {
	if _, err := os.Stdout.WriteString(str); err != nil {
		return err
	}
	return nil
}

func (c ConsoleWriter) Flush() {
	_ = os.Stdout.Sync()
}

type FileWriter struct {
	name      string
	file      *os.File
	buf       *[]byte
	bufSize   int
	lastWrite *time.Time
}

func NewFileWriter(filename string) *FileWriter {
	return &FileWriter{
		bufSize:   1024,
		name:      filename,
		buf:       new([]byte),
		lastWrite: new(time.Time),
	}
}

func (f FileWriter) SetBufferSize(size int) {
	f.bufSize = size
}

func (f FileWriter) Write(str string) error {
	defer func() { *f.lastWrite = time.Now() }()
	str = str + "\n"

	if f.file == nil {
		fl, err := os.OpenFile(f.name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to open file '%s': %s", f.name, err.Error()))
		}
		f.file = fl
	}

	if time.Now().Add(-(time.Second * 10)).After(*f.lastWrite) {
		if err := f.writeBuf(); err != nil {
			return err
		}
		if err := f.writeStr(str); err != nil {
			return err
		}
	}
	if len(str)+len(*f.buf) >= f.bufSize && len(*f.buf) > 0 {
		if err := f.writeBuf(); err != nil {
			return err
		}
	}
	if len(str) >= f.bufSize {
		if err := f.writeBuf(); err != nil {
			return err
		}
		if err := f.writeStr(str); err != nil {
			return err
		}
	}
	*f.buf = append(*f.buf, str...)

	return nil
}

func (f FileWriter) Flush() {
	_ = f.Write("")
	_ = f.file.Sync()
}

func (f FileWriter) writeBuf() error {
	if _, err := f.file.Write(*f.buf); err != nil {
		return errors.New(fmt.Sprintf("failed to write to file '%s': %s", f.name, err))
	}
	*f.buf = (*f.buf)[:0]
	return nil
}

func (f FileWriter) writeStr(str string) error {
	if _, err := f.file.WriteString(str); err != nil {
		return errors.New(fmt.Sprintf("failed to write to file '%s': %s", f.name, err))
	}
	return nil
}
