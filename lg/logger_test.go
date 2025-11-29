package lg

import (
	"os"
	"testing"
)

type wri struct {
	file *os.File
}

func (w wri) Write(message string) error {
	if w.file == nil {
		w.file, _ = os.Open(os.DevNull)
	}
	_, _ = w.file.Write([]byte(message))
	return nil
}

func BenchmarkLogger(b *testing.B) {
	SetCustomPipes(NewPipe(ConsoleSerializer{DisableColors: false}, wri{}))

	b.Run("SingleThreadSync", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			LogSync(LogLevelInfo, "test")
		}
	})

	b.Run("SingleThreadAsync", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Log(LogLevelInfo, "test")
		}
	})

	b.Run("MultiThreadSync", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				LogSync(LogLevelInfo, "test")
			}
		})
	})

	b.Run("MultiThreadAsync", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Log(LogLevelInfo, "test")
			}
		})
	})
}
