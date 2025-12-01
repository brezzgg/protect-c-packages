package lg

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

type wri struct {
	file *os.File
}

func (w *wri) Write(message string) error {
	if w.file == nil {
		w.file, _ = os.Open(os.DevNull)
	}
	_, _ = w.file.WriteString(message)
	return nil
}

func (w *wri) Flush() {}

func BenchmarkLogger(b *testing.B) {
	f, _ := os.CreateTemp(os.TempDir(), "benchmarkLogger-*")
	_ = f.Close()

	fmt.Printf("file: %s\n", f.Name())

	bench := []struct {
		name string
		pipe *Pipe
		msg  string
		args C
	}{
		{
			name: "ConsoleSerializer",
			pipe: NewPipe(ConsoleSerializer{DisableColors: false}, &wri{}),
			msg:  "some test text",
			args: C{"int": 1, "str": "str", "err": errors.New("err")},
		},
		{
			name: "JsonSerializer",
			pipe: NewPipe(JsonSerializer{}, &wri{}),
			msg:  "some test text",
			args: C{"int": 1, "str": "str", "err": errors.New("err")},
		},
		{
			name: "FileWriter-ConsoleSerializer",
			pipe: NewPipe(ConsoleSerializer{DisableColors: false}, NewFileWriter(f.Name())),
			msg:  "some test text",
			args: C{"int": 1, "str": "str", "err": errors.New("err")},
		},
		{
			name: "FileWriter-JsonSerializer",
			pipe: NewPipe(JsonSerializer{}, NewFileWriter(f.Name())),
			msg:  "some test text",
			args: C{"int": 1, "str": "str", "err": errors.New("err")},
		},
	}

	for _, bb := range bench {
		SetCustomPipes(bb.pipe)

		b.Run(bb.name+"-SyngleThreadAsync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Log(LogLevelInfo, bb.msg, bb.args)
			}
		})

		b.Run(bb.name+"-MultiThreadAsync", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					Log(LogLevelInfo, bb.msg, bb.args)
				}
			})
		})
	}

	Close()
	_ = os.Remove(f.Name())
}
