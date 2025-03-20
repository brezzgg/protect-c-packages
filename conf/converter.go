package conf

import (
	"fmt"
	"os"
	"strconv"
)

type Converter struct {
	raw string
}

func NewConverter(raw string) Converter {
	return Converter{raw: raw}
}

func (c Converter) fatal(t string, err error) {
	fmt.Printf("converter: falied to parse value '%s' to %s: %s", c.raw, t, err.Error())
	os.Exit(1)
}

func (c Converter) Int() (int, error) {
	val, err := strconv.Atoi(c.raw)
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c Converter) MustInt() (val int) {
	val, err := c.Int()
	if err != nil {
		c.fatal("int", err)
	}
	return
}

func (c Converter) Float() (float64, error) {
	val, err := strconv.ParseFloat(c.raw, 64)
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c Converter) MustFloat() (val float64) {
	val, err := c.Float()
	if err != nil {
		c.fatal("float", err)
	}
	return
}

func (c Converter) Bool() (bool, error) {
	val, err := strconv.ParseBool(c.raw)
	if err != nil {
		return false, err
	}
	return val, nil
}

func (c Converter) MustBool() (val bool) {
	val, err := c.Bool()
	if err != nil {
		c.fatal("bool", err)
	}
	return
}

func (c Converter) Uint() (uint64, error) {
	val, err := strconv.ParseUint(c.raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c Converter) MustUint() (val uint64) {
	val, err := c.Uint()
	if err != nil {
		c.fatal("uint", err)
	}
	return
}

func (c Converter) Complex() (complex128, error) {
	val, err := strconv.ParseComplex(c.raw, 128)
	if err != nil {
		return -1, err
	}
	return val, nil
}

func (c Converter) MustComplex() (val complex128) {
	val, err := c.Complex()
	if err != nil {
		c.fatal("complex", err)
	}
	return
}
