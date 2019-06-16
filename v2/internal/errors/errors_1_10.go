// +build go1.10

package errors

import (
	"strings"
)

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
