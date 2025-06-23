package utils

import (
	"fmt"
	"runtime"
)

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}
	return fmt.Errorf("error at %s:%d - %v", file, line, err)
}
