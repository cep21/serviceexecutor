package serviceexecutor

import "strings"

type multiError struct {
	errs []error
}

func (m *multiError) Error() string {
	ret := make([]string, len(m.errs))
	for i, err := range m.errs {
		ret[i] = err.Error()
	}
	return strings.Join(ret, ",")
}

func errFromManyErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	nonNil := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNil = append(nonNil, err)
		}
	}
	if len(nonNil) == 0 {
		return nil
	}
	if len(nonNil) == 1 {
		return nonNil[0]
	}
	return &multiError{errs: nonNil}
}
