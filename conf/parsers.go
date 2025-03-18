package conf

import (
	"flag"
	"os"
)

type ArgParser struct {
	params []Parameter
}

type linkParam struct {
	link  *string
	param Parameter
}

func (a ArgParser) Parse(keyConvertFunc func(k string) string) (r map[string]string) {
	temp := make(map[string]linkParam)
	for _, param := range a.params {
		temp[param.Key()] = linkParam{
			link:  flag.String(param.Key(), "", ""),
			param: param,
		}
	}
	flag.Parse()

	r = make(map[string]string)
	for key, linkPar := range temp {
		resultKey := keyConvertFunc(key)
		resultValue := linkPar.param.String(*linkPar.link)

		if linkPar.param.Dest() != nil {
			*linkPar.param.Dest() = resultValue
		}
		if linkPar.param.DestConv() != nil {
			*linkPar.param.DestConv() = NewConverter(resultValue)
		}

		r[resultKey] = resultValue
	}
	return
}

func Arg(params ...Parameter) ArgParser {
	return ArgParser{params: params}
}

type EnvParser struct {
	params []Parameter
}

func (e EnvParser) Parse(keyConvertFunc func(k string) string) (r map[string]string) {
	r = make(map[string]string)
	for _, param := range e.params {
		resultKey := keyConvertFunc(param.Key())
		resultValue := param.String(os.Getenv(param.Key()))

		if param.Dest() != nil {
			*param.Dest() = resultValue
		}
		if param.DestConv() != nil {
			*param.DestConv() = NewConverter(resultValue)
		}

		r[resultKey] = resultValue
	}
	return
}

func Env(params ...Parameter) EnvParser {
	return EnvParser{params: params}
}
