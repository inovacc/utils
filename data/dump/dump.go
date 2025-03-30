package dump

import (
	"github.com/inovacc/gopickle"
)

func Dump(filename string, data any) error {
	return gopickle.Dump(filename, data)
}

func Load(filename string, data any) error {
	return gopickle.Load(filename, data)
}
