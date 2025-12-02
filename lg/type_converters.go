package lg

type defaultTypeConverter struct{}

func NewDefaultTypeConverter() *defaultTypeConverter {
	return &defaultTypeConverter{}
}

func (c defaultTypeConverter) ConvAndPush(item any, push func(key string, val any)) {
	switch value := item.(type) {
	case C:
		c.convAndPushContext(value, push)
	case error:
		push("error", value.Error())
	default:
		push("arg", value)
	}
}

func (c defaultTypeConverter) convAndPushContext(context C, push func(key string, val any)) {
	for key, val := range context {
		switch v := val.(type) {
		case error:
			push(key, v.Error())
		default:
			push(key, v)
		}
	}
}
