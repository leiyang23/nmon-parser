package util

import (
	"fmt"
	"io"
)

func Close(closer io.Closer) {
	_ = closer.Close()
}

func WrapErr(err error, s string) error {
	return fmt.Errorf("%s: %v", s, err)
}
