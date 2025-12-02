package lg

/*
C is an alias that is equivalent to map[string]any.
It is used to add context to logs in a convenient way.
*/
type C map[string]any

type TypeConverter interface {
	ConvAndPush(item any, push func(key string, val any))
}

type Serializer interface {
	Serialize(m Message) string
}

type Writer interface {
	Write(Message string) error
	Flush()
	Close()
}

type ColorizedOut interface {
	Colorize(color string) string
	Colorized() string
	Standard() string
}
