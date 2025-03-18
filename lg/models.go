package lg

var (
	LogLevelDebug       = NewLogLevel("Debug", ClrFgBoldGreen)
	LogLevelInfo        = NewLogLevel("Info", ClrFgBoldBlue)
	LogLevelWarn        = NewLogLevel("Warn", ClrFgBoldYellow)
	LogLevelError       = NewLogLevel("Error", ClrFgBoldRed)
	LogLevelFatal       = NewLogLevel("Fatal", ClrFgBoldRed)
	logLevelLoggerError = NewLogLevel("LoggerError", ClrFgBoldPink)
)

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

/*
C is an alias that is equivalent to map[string]any.
It is used to add context to logs in a convenient way.
*/
type C map[string]any

/*
E is an alias that is equivalent to struct{Err error}.
This is used to explicitly specify the error in the context of the log,
and to make sure that if the value is nil,
the log will display {"error": "nil"}.
*/
type E struct {
	Err error
}

type TypeConverter interface {
	ConvAndPush(item any, push func(key string, val any))
}

type Serializer interface {
	Serialize(m Message) string
}

type Writer interface {
	Write(Message string) error
}

type ColorizedOut interface {
	Colorize(color string) string
	Colorized() string
	Standard() string
}
