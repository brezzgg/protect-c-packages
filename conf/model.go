package conf

type Parameter interface {
	Key() string
	String(val string) string
	Write(dest *string) Parameter
	Dest() *string
	WriteConv(dest *Converter) Parameter
	DestConv() *Converter
}

type Parser interface {
	Parse(keyConvertFunc func(k string) string) map[string]string
}
