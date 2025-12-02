package lg

import (
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type (
	consoleSerializer struct {
		disableColors bool
		cachedOffset  int
		cachedTime    time.Time
	}

	ConsoleSerializerOption func(*consoleSerializer)
)

func NewConsoleSerializer(opts ...ConsoleSerializerOption) *consoleSerializer {
	s := &consoleSerializer{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithDisabledColors() ConsoleSerializerOption {
	return func(s *consoleSerializer) {
		s.disableColors = true
	}
}

func (c *consoleSerializer) Serialize(m Message) string {
	context := ""
	if m.Context != nil {
		b, err := json.Marshal(m.Context)
		str := string(b)
		if err == nil && len(str) > 0 && str != "{}" {
			context = str
		}
	}

	offset := c.getTimeOffset(time.Now())
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
	if c.disableColors {
		caller = m.Caller.Standard()
		level = m.Level.Standard() + levelSpaces
	} else {
		caller = m.Caller.Colorize(ClrFgCyan)
		level = m.Level.Colorized() + levelSpaces
	}

	return m.Time.UTC().Format("2006/01/02 15:04:05") + offsetStr + "  " + level + "  " + caller + "  " + m.Text + " " + context
}

func (c *consoleSerializer) getTimeOffset(t time.Time) int {
	if c.cachedTime.YearDay() != t.YearDay() {
		_, c.cachedOffset = t.Zone()
		c.cachedTime = t
	}
	return c.cachedOffset
}
