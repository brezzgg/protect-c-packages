package lg

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LogLevel struct {
	Levels []string
	color  string
}

func NewLogLevel(asciiClr string, levels ...string) LogLevel {
	for i := range levels {
		levels[i] = strings.TrimSpace(levels[i])
	}
	return LogLevel{
		Levels: levels,
		color:  asciiClr,
	}
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Levels)
}

func (l LogLevel) Colorize(color string) string {
	return fmt.Sprintf("%s%s%s", color, l.Formatter(), ClrReset)
}

func (l LogLevel) Colorized() string {
	return fmt.Sprintf("%s%s%s", l.color, l.Formatter(), ClrReset)
}

func (l LogLevel) Standard() string {
	return l.Formatter()
}

func (l LogLevel) Formatter() string {
	return strings.Join(l.Levels, ".")
}

type Caller struct {
	Method string `json:"method"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	color  string
}

func (c Caller) Colorize(color string) string {
	return fmt.Sprintf("%s%s%s", color, c.formatter(), ClrReset)
}

func (c Caller) Colorized() string {
	return fmt.Sprintf("%s%s%s", c.color, c.formatter(), ClrReset)
}

func (c Caller) Standard() string {
	return c.File
}

func (c Caller) formatter() string {
	return fmt.Sprintf("%s.%s:%d", c.File, c.Method, c.Line)
}
