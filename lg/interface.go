package lg

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

var (
	GlobalLogger = NewLogger()
)

// --------------- Logger struct functions ---------------

/*
SetCustomPipes allows you to customize the global logger.
*/
func (l *Logger) SetCustomPipes(pipes ...*Pipe) {
	l.pipes = pipes
}

/*
Debug function outputs a log with the level LogLevelDebug.
*/
func (l *Logger) Debug(args ...any) {
	l.Handle(args, LogLevelDebug)
}

/*
Info function outputs a log with the level LogLevelInfo.
*/
func (l *Logger) Info(args ...any) {
	l.Handle(args, LogLevelInfo)
}

/*
Warn function outputs a log with the level LogLevelWarn.
*/
func (l *Logger) Warn(args ...any) {
	l.Handle(args, LogLevelWarn)
}

/*
Error function outputs a log with the level LogLevelError.
*/
func (l *Logger) Error(args ...any) {
	l.Handle(args, LogLevelError)
}

/*
Fatal function outputs a log with the level LogLevelFatal
and terminates the program with exit code 1,
equivalent of os.Exit(1).
*/
func (l *Logger) Fatal(args ...any) {
	// execute end tasks
	endTasks.Execute()

	l.Handle(args, LogLevelFatal)

	l.Close()

	// exit with code 1
	os.Exit(1)
}

/*
Panic function outputs a log with the level LogLevelPanic,
terminates the program with exit code 1 and print stack trace.
*/
func (l *Logger) Panic(args ...any) {
	// execute end tasks
	endTasks.Execute()

	l.Handle(args, LogLevelPanic)

	l.Close()

	// print stack trace and exit with code 1
	debug.PrintStack()
	os.Exit(1)
}

/*
Exit function terminate the program with your exit code
*/
func (l *Logger) Exit(code int) {
	// execute end tasks
	endTasks.Execute()

	l.Close()

	// exit
	os.Exit(code)
}

/*
Log function outputs a log with your custom LogLevel.
*/
func (l *Logger) Log(level LogLevel, args ...any) {
	l.Handle(args, level)
}

/*
Invoked function outputs a log with the level LogLevelDebug,
which notifies you that a function has been called.
*/
func (l *Logger) Invoked() {
	l.Handle([]any{"Invoked"}, LogLevelDebug)
}

/*
End function returns structure instance EndTasks.
*/
func (l *Logger) End() *EndTasks {
	return endTasks
}

// --------------- Global logger functions ---------------

/*
Debug function invokes the equivalent function Logger.Debug for the global logger.
*/
func Debug(args ...any) {
	GlobalLogger.Debug(args...)
}

/*
Info function invokes the equivalent function Logger.Info for the global logger.
*/
func Info(args ...any) {
	GlobalLogger.Info(args...)
}

/*
Warn function invokes the equivalent function Logger.Warn for the global logger.
*/
func Warn(args ...any) {
	GlobalLogger.Warn(args...)
}

/*
Error function invokes the equivalent function Logger.Error for the global logger.
*/
func Error(args ...any) {
	GlobalLogger.Error(args...)
}

/*
Fatal function invokes the equivalent function Logger.Fatal for the global logger.
*/
func Fatal(args ...any) {
	GlobalLogger.Fatal(args...)
}

/*
Panic function invokes the equivalent function Logger.Panic for the global logger.
*/
func Panic(args ...any) {
	GlobalLogger.Panic(args...)
}

/*
Log function invokes the equivalent function Logger.Log for the global logger.
*/
func Log(level LogLevel, args ...any) {
	GlobalLogger.Log(level, args...)
}

/*
Close function invokes the equivalent function Logger.Close for the global logger.
*/
func Close() {
	GlobalLogger.Close()
}

/*
F function is alias for fmt.Sprintf function.
*/
func F(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

/*
T function is alias for strings.TrimSpaces function.
*/
func T(str string) string {
	return strings.TrimSpace(str)
}
