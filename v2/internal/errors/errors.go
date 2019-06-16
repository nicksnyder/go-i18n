package errors

type multierr []error

func Append(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	me1, ok1 := err1.(multierr)
	me2, ok2 := err2.(multierr)
	switch {
	case ok1 && !ok2:
		return append(me1, err2)
	case ok1 && ok2:
		return append(me1, me2...)
	case !ok1 && ok2:
		return append(multierr{err1}, me2...)
	case !ok1 && !ok2:
		return multierr{err1, err2}
	}
	panic("unreachable")
}
