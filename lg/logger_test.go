package lg

import (
	"errors"
	"os"
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	null, _ := os.Open(os.DevNull)
	defer null.Close()

	tmp, _ := os.CreateTemp(os.TempDir(), "benchmarkLogger-*")
	_ = tmp.Close()

	args := C{
		"int":    9812738949823,
		"string": "some string",
		"err":    errors.New("some error"),
		"struct": struct {
			arg1 string
			arg2 string
			arg3 string
			arg4 string
			arg5 string
		}{
			"a",
			"b",
			"c",
			"d",
			"e",
		},
	}

	bench := []struct {
		name  string
		pipe  *Pipe
		pipes []*Pipe
		msg   string
		args  C
	}{
		{
			name: "consoleSerializer",
			pipe: NewPipe(WithSerializer(NewConsoleSerializer()), WithWriter(NewConsoleWriter(WithCustomStdout(null)))),
			msg:  "some test text",
			args: args,
		},
		{
			name: "JsonSerializer",
			pipe: NewPipe(WithSerializer(NewJSONSerializer()), WithWriter(NewConsoleWriter(WithCustomStdout(null)))),
			msg:  "some test text",
			args: args,
		},
		{
			name: "FileWriter-consoleSerializer",
			pipe: NewPipe(WithSerializer(NewConsoleSerializer()), WithWriter(NewFileWriter(tmp.Name()))),
			msg:  "some test text",
			args: args,
		},
		{
			name: "FileWriter-JsonSerializer",
			pipe: NewPipe(WithSerializer(NewJSONSerializer()), WithWriter(NewFileWriter(tmp.Name()))),
			msg:  "some test text",
			args: args,
		},
		{
			pipes: []*Pipe{
				NewPipe(WithSerializer(NewConsoleSerializer()), WithWriter(NewConsoleWriter(WithCustomStdout(null)))),
				NewPipe(WithSerializer(NewJSONSerializer()), WithWriter(NewFileWriter(tmp.Name()))),
			},
			name: "MultiPipe",
			msg:  "some test text",
			args: args,
		},
	}

	for _, bb := range bench {
		var logger *Logger
		if bb.pipe != nil {
			logger = NewLogger(WithPipe(bb.pipe))
		} else {
			logger = NewLogger(WithPipes(bb.pipes...))
		}

		b.Run(bb.name+"-SyngleThreadAsync", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logger.Log(LogLevelInfo, bb.msg, bb.args)
			}
		})

		b.Run(bb.name+"-MultiThreadAsync", func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					logger.Log(LogLevelInfo, bb.msg, bb.args)
				}
			})
		})

		logger.Close()
	}
	_ = os.Remove(tmp.Name())
}
