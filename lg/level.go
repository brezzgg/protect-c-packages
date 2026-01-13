package lg

var (
	LogLevelDebug       = NewLogLevel(ClrFgBoldGreen, "Debug").WithPriority(100)
	LogLevelInfo        = NewLogLevel(ClrFgBoldBlue, "Info").WithPriority(200)
	LogLevelWarn        = NewLogLevel(ClrFgBoldYellow, "Warn").WithPriority(300)
	LogLevelError       = NewLogLevel(ClrFgBoldRed, "Error").WithPriority(400)
	LogLevelFatal       = NewLogLevel(ClrFgBoldRed, "Fatal").WithPriority(500)
	LogLevelPanic       = NewLogLevel(ClrFgBoldRed, "Panic").WithPriority(600)
	logLevelLoggerError = NewLogLevel(ClrFgBoldPink, "LoggerError").WithPriority(65535)
)

const (
	LevelOptionDisableCaller     LogLevelOption = "caller_disable"
	LevelOptionCallerOnlyFunc    LogLevelOption = "caller_only_func"
	LevelOptionCallerOnlyFile    LogLevelOption = "caller_only_file"
	LevelOptionCallerDisableFunc LogLevelOption = "caller_disable_func"
	LevelOptionCallerDisableFile LogLevelOption = "caller_disable_file"
	LevelOptionCallerDisableLine LogLevelOption = "caller_disable_line"

	LevelOptionDisableTime       LogLevelOption = "time_disable"
	LevelOptionTimeDisableOffset LogLevelOption = "time_disable_offset"
)

const LevelDefaultPriority = uint16(200)

const (
	ClrReset = "\033[0;0m"

	ClrFgRed    = "\033[0;91m"
	ClrFgGreen  = "\033[0;92m"
	ClrFgYellow = "\033[0;93m"
	ClrFgBlue   = "\033[0;94m"
	ClrFgPink   = "\033[0;95m"
	ClrFgCyan   = "\033[0;96m"
	ClrFgGray   = "\033[0;97m"
	ClrFgWhite  = "\033[0;98m"

	ClrFgBoldRed    = "\033[1;91m"
	ClrFgBoldGreen  = "\033[1;92m"
	ClrFgBoldYellow = "\033[1;93m"
	ClrFgBoldBlue   = "\033[1;94m"
	ClrFgBoldPink   = "\033[1;95m"
	ClrFgBoldCyan   = "\033[1;96m"
	ClrFgBoldGray   = "\033[1;97m"
	ClrFgBoldWhite  = "\033[1;98m"
)
