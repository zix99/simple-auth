package selector

import "fmt"

func jsonErrorf(code string, underlyingErrors []error, s string, args ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"error":   true,
		"message": fmt.Sprintf(s, args...),
		"reason":  code,
		"cause":   underlyingErrors,
	}
}
