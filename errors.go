package bipf

import "fmt"

func wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%v: %w", message, err)
}

func wrapf(err error, format string, args ...any) error {
	return wrap(err, fmt.Sprintf(format, args...))
}
