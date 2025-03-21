package lg

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type ConsoleSerializer struct {
	DisableColors bool
}

func (c ConsoleSerializer) Serialize(m Message) string {
	context := ""
	if m.Context != nil {
		b, err := json.Marshal(m.Context)
		str := string(b)
		if err == nil && len(str) > 0 && str != "{}" {
			context = str
		}
	}

	_, offset := m.Time.Zone()
	offsetStr := strconv.Itoa(offset / 3600)
	if !strings.HasPrefix(offsetStr, "-") {
		offsetStr = "+" + offsetStr
	}

	var caller, level string
	levelSpaces := ""

	if form := m.Level.Formatter(); len(form) < 5 {
		for len(form)+len(levelSpaces) < 5 {
			levelSpaces += " "
		}
	}
	if c.DisableColors {
		caller = m.Caller.Standard()
		level = m.Level.Standard() + levelSpaces
	} else {
		caller = m.Caller.Colorize(ClrFgCyan)
		level = m.Level.Colorized() + levelSpaces
	}

	return fmt.Sprintf("%s  %s  %s  %s %s",
		m.Time.UTC().Format("2006/01/02 15:04:05")+offsetStr,
		level,
		caller,
		m.Text,
		context,
	)
}

type JsonSerializer struct{}

func (j JsonSerializer) Serialize(m Message) string {
	b, _ := json.Marshal(m)
	return string(b)
}
