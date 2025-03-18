package lg

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LogLevel struct {
	Name  string
	color string
}

func NewLogLevel(name string, asciiColor string) LogLevel {
	return LogLevel{name, asciiColor}
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(l.Name))
}

func (l LogLevel) Colorize(color string) string {
	return fmt.Sprintf("%s%s%s", color, l.Name, ClrReset)
}

func (l LogLevel) Colorized() string {
	return fmt.Sprintf("%s%s%s", l.color, l.Name, ClrReset)
}

func (l LogLevel) Standard() string {
	return l.Name
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
