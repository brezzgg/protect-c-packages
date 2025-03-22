package lg

const FormatedErrorKey = "$formattedError"

type DefaultTypeConverter struct{}

func (DefaultTypeConverter) ConvAndPush(item any, push func(key string, val any)) {
	if ctx, ok := item.(C); ok {
		for key, value := range ctx {
			push(key, value)
		}
	} else if err, ok := item.(error); ok {
		push("error", err.Error())
	} else if errS, ok := item.(E); ok {
		v := ""
		if errS.Err == nil {
			v = "nil"
		} else {
			v = errS.Err.Error()
		}
		push("error", v)
	} else if str, ok := item.(string); ok {
		push("arg", str)
	}
}
