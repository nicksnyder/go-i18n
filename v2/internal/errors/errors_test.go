package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestAppend(t *testing.T) {
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
			err1:        errors.New("err1"),
			err2:        nil,
			expectedErr: errors.New("err1"),
		},
		{
			err1:        nil,
			err2:        errors.New("err2"),
			expectedErr: errors.New("err2"),
		},
		{
			err1:        errors.New("err1"),
			err2:        errors.New("err2"),
			expectedErr: multierr{errors.New("err1"), errors.New("err2")},
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
