package conf

import (
	"errors"
	"fmt"
)

type ParameterWithDefault struct {
	Default, key string
	dest         *string
	destConv     *Converter
}

func (p ParameterWithDefault) String(val string) string {
	if len(val) == 0 {
		return p.Default
	} else {
		return val
	}
}

func (p ParameterWithDefault) Key() string {
	return p.key
}

func (p ParameterWithDefault) Write(dest *string) Parameter {
	p.dest = dest
	return p
}

func (p ParameterWithDefault) Dest() *string {
	return p.dest
}

func (p ParameterWithDefault) WriteConv(dest *Converter) Parameter {
	p.destConv = dest
	return p
}

func (p ParameterWithDefault) DestConv() *Converter {
	return p.destConv
}

func DParam(key, def string) ParameterWithDefault {
	return ParameterWithDefault{key: key, Default: def}
}

type ParameterRequired struct {
	key      string
	dest     *string
	destConv *Converter
}

func (p ParameterRequired) String(val string) string {
	if len(val) == 0 {
		parseError = errors.New(fmt.Sprintf("required parameter '%s' is missing", p.key))
		return ""
	}
	return val
}

func (p ParameterRequired) Key() string {
	return p.key
}

func (p ParameterRequired) Write(dest *string) Parameter {
	p.dest = dest
	return p
}

func (p ParameterRequired) Dest() *string {
	return p.dest
}

func (p ParameterRequired) WriteConv(dest *Converter) Parameter {
	p.destConv = dest
	return p
}

func (p ParameterRequired) DestConv() *Converter {
	return p.destConv
}

func RParam(key string) ParameterRequired {
	return ParameterRequired{key: key}
}
