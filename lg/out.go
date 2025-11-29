package lg

import (
	"slices"
	"strconv"
	"strings"

	"encoding/json"
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
	var sb strings.Builder
	sb.Grow(len(color) + len(ClrReset) + 50)
	sb.WriteString(color)
	sb.WriteString(l.Formatter())
	sb.WriteString(ClrReset)
	return sb.String()
}

func (l LogLevel) Colorized() string {
	var sb strings.Builder
	sb.Grow(len(l.color) + len(ClrReset) + 50)
	sb.WriteString(l.color)
	sb.WriteString(l.Formatter())
	sb.WriteString(ClrReset)
	return sb.String()
}

func (l LogLevel) Standard() string {
	return l.Formatter()
}

func (l LogLevel) AppendLevels(levels ...string) LogLevel {
	l.Levels = append(l.Levels, levels...)
	return l
}

func (l LogLevel) Formatter() string {
	return strings.Join(l.Levels, ".")
}

func (l LogLevel) Equal(other LogLevel) bool {
	return l.color == other.color && slices.Equal(l.Levels, other.Levels)
}

type Caller struct {
	Method string `json:"method"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	color  string
}

func (c Caller) Colorize(color string) string {
	var sb strings.Builder
	sb.Grow(len(color) + len(ClrReset) + 50)
	sb.WriteString(color)
	sb.WriteString(c.formatter())
	sb.WriteString(ClrReset)
	return sb.String()
}

func (c Caller) Colorized() string {
	var sb strings.Builder
	sb.Grow(len(c.color) + len(ClrReset) + 50)
	sb.WriteString(c.color)
	sb.WriteString(c.formatter())
	sb.WriteString(ClrReset)
	return sb.String()
}

func (c Caller) Standard() string {
	return c.File
}

func (c Caller) formatter() string {
	var sb strings.Builder
	sb.Grow(len(c.File) + len(c.Method))
	sb.WriteString(c.File)
	sb.WriteString(".")
	sb.WriteString(c.Method)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(c.Line))
	return sb.String()
}

func (c Caller) Equal(other Caller) bool {
	return c.Line == other.Line &&
		c.Method == other.Method &&
		c.File == other.File &&
		c.color == other.color
}
