package errors

type multierr []error

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
