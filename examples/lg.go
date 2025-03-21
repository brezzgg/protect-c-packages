package main

import "github.com/brezzgg/protect-c-packages/lg"

func main() {
	ExampleSetupLogger()
	ExampleLogLevels()
}

func ExampleSetupLogger() {
	lg.SetCustomPipes(
		lg.NewPipe(
			lg.ConsoleSerializer{
				DisableColors: false,
			},
			lg.NewConsoleWriter(),
		),
	)

	lg.Info("some log", lg.C{"with some": "context"})
	// output in console: 2025/01/01 00:00:00+0  Info   main.ExampleSetupLogger:29  some log {"with some":"context"}

	lg.SetCustomPipes(
		lg.NewPipe(
			lg.JsonSerializer{},
			lg.NewFileWriter("log.json"),
		),
	)

	lg.Info("some log", lg.C{"with some": "context"})
	// output in log.json: {"time":"2025-01-01T00:00:00.0000000+00:00","caller":{"method":"ExampleSetupLogger","file":"main","line":20},"level":"info","msg":"some log","ctx":{"with some":"context"}}
}

func ExampleLogLevels() {
	lg.SetCustomPipes(lg.DefaultConsolePipe(false))

	var (
		str = "Hello world"
		ctx = lg.C{"hello": "world"}
	)

	// create custom log level
	logLevelProd := lg.NewLogLevel("Prod", lg.ClrFgCyan)

	// write log with custom log level
	lg.Log(logLevelProd, str, ctx)

	lg.Debug(str, ctx)
	lg.Info(str, ctx)
	lg.Warn(str, ctx)
	lg.Error(str, ctx)
	lg.Fatal(str, ctx)
}
