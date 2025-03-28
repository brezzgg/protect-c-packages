package lg

import (
	"encoding/json"
	"errors"
	"os"
	"runtime/debug"
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
	args = errorDelivery(args...)
	globalLogger.Handle(args, LogLevelError, nil)
}

/*
Fatal function outputs a log with the level LogLevelFatal
and terminates the program with exit code 1,
equivalent of os.Exit(1).
*/
func Fatal(args ...any) {
	// execute end tasks
	globalLogger.EndTasks.Execute()

	args = errorDelivery(args...)

	readyCh := make(chan any)
	globalLogger.Handle(args, LogLevelFatal, readyCh)

	<-readyCh
	// exit with code 1
	os.Exit(1)
}

/*
Panic function outputs a log with the level LogLevelPanic,
terminates the program with exit code 1 and print stack trace.
*/
func Panic(args ...any) {
	// execute end tasks
	globalLogger.EndTasks.Execute()

	args = errorDelivery(args...)

	readyCh := make(chan any)
	globalLogger.Handle(args, LogLevelPanic, readyCh)

	<-readyCh
	// print stack trace and exit with code 1
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
LogSync is equivalent of Log function, but it runs in sync.
*/
func LogSync(level LogLevel, args ...any) {
	readyCh := make(chan any)
	globalLogger.Handle(args, level, readyCh)

	<-readyCh
}

/*
FormatError is a function that creates an error recognizable by the Logger,
which, when passed to the Fatal or Error function, will be displayed as,
as it was created.
*/
func FormatError(args ...any) error {
	var msg Message
	err := globalLogger.CheckArgs(args, &msg.Text, &msg.Context)
	if err != nil {
		Log(logLevelLoggerError, "Failed to create formatted error", C{"args": args})
	}
	msg.Caller = GetCallerInfo(1)

	b, err := json.Marshal(msg)
	if err != nil {
		Log(logLevelLoggerError, "Failed to marshal formatted error", C{"args": args})
	}
	return errors.New(FormatedErrorKey + string(b))
}

/*
Invoked function outputs a log with the level LogLevelDebug,
which notifies you that a function has been called.
*/
func Invoked() {
	globalLogger.Handle([]any{"Invoked"}, LogLevelDebug, nil)
}

/*
End function returns structure instance EndTasks.
*/
func End() *EndTasks {
	return globalLogger.EndTasks
}

func errorDelivery(args ...any) []any {
	if err, ok := args[0].(error); ok {
		args[0] = "Undescribed error: "
		args = append(args, err)
	}
	return args
}
