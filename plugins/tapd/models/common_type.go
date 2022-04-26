package models

import (
	"strconv"
	"strings"
)

type Uint64s uint64

func (u *Uint64s) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		*u = Uint64s(0)
		return nil
	}

	str = strings.Trim(str, `"`)
	if str == "-1" {
		*u = Uint64s(0)
		return nil
	}
	ui, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	*u = Uint64s(ui)
	return nil
}

type Ints int

func (i *Ints) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		*i = Ints(0)
		return nil
	}
	str = strings.Trim(str, `"`)
	ui, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return err
	}
	*i = Ints(ui)
	return nil
}
