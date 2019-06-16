package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestAppend(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")
	err4 := errors.New("err4")
	tests := []struct {
		err1        error
		err2        error
		expectedErr error
	}{
		{
			err1:        nil,
			err2:        nil,
			expectedErr: nil,
		},
		{
			err1:        err1,
			err2:        nil,
			expectedErr: err1,
		},
		{
			err1:        nil,
			err2:        err2,
			expectedErr: err2,
		},
		{
			err1:        err1,
			err2:        err2,
			expectedErr: multierr{err1, err2},
		},
		{
			err1:        multierr{err1},
			err2:        err2,
			expectedErr: multierr{err1, err2},
		},
		{
			err1:        multierr{err1, err2},
			err2:        err3,
			expectedErr: multierr{err1, err2, err3},
		},
		{
			err1:        err1,
			err2:        multierr{err2},
			expectedErr: multierr{err1, err2},
		},
		{
			err1:        err1,
			err2:        multierr{err2, err3},
			expectedErr: multierr{err1, err2, err3},
		},
		{
			err1:        multierr{err1, err2},
			err2:        multierr{err3, err4},
			expectedErr: multierr{err1, err2, err3, err4},
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("Append(%q,%q)=%q", test.err1, test.err2, test.expectedErr)
		t.Run(name, func(t *testing.T) {
			actual := Append(test.err1, test.err2)
			if !reflect.DeepEqual(actual, test.expectedErr) {
				t.Fatalf("expected %q; got %q", test.expectedErr, actual)
			}
		})
	}
}
