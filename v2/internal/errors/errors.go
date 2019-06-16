package errors

import (
	"strings"
)

type multierr []error

func (e multierr) Error() string {
	var b strings.Builder
	for i, err := range e {
		if i != 0 {
			b.WriteString("; ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

func Append(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	me, ok := err1.(multierr)
	if !ok {
		return multierr{err1, err2}
	}
	return append(me, err2)
}
