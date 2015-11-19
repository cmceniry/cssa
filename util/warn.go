package util

import (
	"fmt"
)

func appendWarn(warns []string, f string, v ...interface{}) []string {
	return append(warns, fmt.Sprintf(f, v...))
}
