package auth

import "fmt"

func jsonErrorf(code, s string, args ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"error":   true,
		"message": fmt.Sprintf(s, args...),
		"reason":  code,
	}
}
