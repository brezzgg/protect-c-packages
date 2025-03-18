package lg

import (
	"os"
	"runtime/debug"
	"time"
)

var (
	globalLogger = DefaultConsoleLogger()
)

/*
SetCustomPipes allows you to customize the global logger.
*/
func SetCustomPipes(pipes ...*Pipe) {
	globalLogger.Pipes = pipes
}

/*
Debug function outputs a log with the level LogLevelDebug.
*/
func Debug(args ...any) {
	globalLogger.Handle(args, LogLevelDebug, nil)
}

/*
Info function outputs a log with the level LogLevelInfo.
*/
func Info(args ...any) {
	globalLogger.Handle(args, LogLevelInfo, nil)
}

/*
Warn function outputs a log with the level LogLevelWarn.
*/
func Warn(args ...any) {
	globalLogger.Handle(args, LogLevelWarn, nil)
}

/*
Error function outputs a log with the level LogLevelError.
*/
func Error(args ...any) {
	globalLogger.Handle(args, LogLevelError, nil)
}

/*
Fatal function outputs a log with the level LogLevelFatal
and terminates the program with code 1,
equivalent of os.Exit(1).
*/
func Fatal(args ...any) {
	readyCh := make(chan any)
	globalLogger.Handle(args, LogLevelFatal, readyCh)

	<-readyCh
	// Print stack trace end exit with code 1
	// To make sure all the logs are output to the console
	time.Sleep(100 * time.Millisecond)
	debug.PrintStack()
	os.Exit(1)
}

/*
Log function outputs a log with your custom LogLevel.
*/
func Log(level LogLevel, args ...any) {
	globalLogger.Handle(args, level, nil)
}

/*
Invoked function outputs a log with the level LogLevelDebug,
which notifies you that a function has been called.
*/
func Invoked() {
	globalLogger.Handle([]any{"Invoked"}, LogLevelDebug, nil)
}
