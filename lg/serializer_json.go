package lg

import (
	"github.com/goccy/go-json"
)

type (
	jsonSerializer struct {
		indent bool
	}

	JSONSerializerOption func(*jsonSerializer)
)

func NewJSONSerializer(opts ...JSONSerializerOption) *jsonSerializer {
	j := &jsonSerializer{}
	for _, opt := range opts {
		opt(j)
	}
	return j
}

func WithIndent() JSONSerializerOption {
	return func(j *jsonSerializer) {
		j.indent = true
	}
}

func (j *jsonSerializer) Serialize(m Message) string {
	var b []byte
	if j.indent {
		b, _ = json.MarshalIndent(m, "", "  ")
	} else {
		b, _ = json.Marshal(m)
	}
	return string(b)
}
