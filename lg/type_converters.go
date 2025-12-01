package lg

type DefaultTypeConverter struct{}

func (c DefaultTypeConverter) ConvAndPush(item any, push func(key string, val any)) {
	switch value := item.(type) {
	case C:
		c.convAndPushContext(value, push)
	case error:
		push("error", value.Error())
	default:
		push("arg", value)
	}
}

func (c DefaultTypeConverter) convAndPushContext(context C, push func(key string, val any)) {
	for key, val := range context {
		switch v := val.(type) {
		case error:
			push(key, v.Error())
		default:
			push(key, v)
		}
	}
}
