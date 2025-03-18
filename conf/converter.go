package conf

import "strconv"

type Converter struct {
	raw string
}

func NewConverter(raw string) Converter {
	return Converter{raw: raw}
}

func (c Converter) Int() int {
	r, err := strconv.Atoi(c.raw)
	if err != nil {
		return -1
	}
	return r
}

func (c Converter) Float() float64 {
	r, err := strconv.ParseFloat(c.raw, 64)
	if err != nil {
		return -1
	}
	return r
}

func (c Converter) Bool() (res bool, ok bool) {
	r, err := strconv.ParseBool(c.raw)
	if err != nil {
		return false, false
	}
	return r, true
}
