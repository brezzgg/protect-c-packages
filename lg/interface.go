package lg

import (
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
	globalLogger.Handle(args, LogLevelDebug)
}

/*
Info function outputs a log with the level LogLevelInfo.
*/
func Info(args ...any) {
	globalLogger.Handle(args, LogLevelInfo)
}

/*
Warn function outputs a log with the level LogLevelWarn.
*/
func Warn(args ...any) {
	globalLogger.Handle(args, LogLevelWarn)
}

/*
Error function outputs a log with the level LogLevelError.
*/
func Error(args ...any) {
	globalLogger.Handle(args, LogLevelError)
}

/*
Fatal function outputs a log with the level LogLevelFatal
and terminates the program with exit code 1,
equivalent of os.Exit(1).
*/
func Fatal(args ...any) {
	// execute end tasks
	globalLogger.EndTasks.Execute()

	globalLogger.Handle(args, LogLevelFatal)

	globalLogger.Close()

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

	globalLogger.Handle(args, LogLevelPanic)

	globalLogger.Close()

	// print stack trace and exit with code 1
	debug.PrintStack()
	os.Exit(1)
}

/*
Exit function terminate the program with your exit code
*/
func Exit(code int) {
	// execute end tasks
	globalLogger.EndTasks.Execute()

	globalLogger.Close()

	// exit
	os.Exit(code)
}

/*
Close function dispose logger
*/
func Close() {
	globalLogger.Close()
}

/*
Log function outputs a log with your custom LogLevel.
*/
func Log(level LogLevel, args ...any) {
	globalLogger.Handle(args, level)
}

/*
Invoked function outputs a log with the level LogLevelDebug,
which notifies you that a function has been called.
*/
func Invoked() {
	globalLogger.Handle([]any{"Invoked"}, LogLevelDebug)
}

/*
End function returns structure instance EndTasks.
*/
func End() *EndTasks {
	return globalLogger.EndTasks
}
