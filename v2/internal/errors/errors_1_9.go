// +build !go1.10

package errors

import "bytes"

func (e multierr) Error() string {
	var b bytes.Buffer
	for i, err := range e {
		if i != 0 {
			b.WriteString("; ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}
