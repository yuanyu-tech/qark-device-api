package device

import (
	"fmt"
	"strconv"
	"strings"
)

func GetInt(t interface{}) int {
	switch v := t.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return -1
		} else {
			return i
		}
	}

	return -1
}

func GetInt64(t interface{}) int64 {
	switch v := t.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return -1
		} else {
			return i
		}
	}

	return -1
}

func GetBool(t interface{}) bool {
	switch v := t.(type) {
	case bool:
		return v
	case string:
		sv := strings.ToLower(v)
		if sv == "true" || sv == "1" || sv == "on" {
			return true
		} else {
			return false
		}
	case int:
		if v > 0 {
			return true
		} else {
			return false
		}
	}

	return false
}

func GetString(t interface{}) string {
	switch v := t.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%d", int(v))
	case int:
		return fmt.Sprintf("%d", v)
	}

	return ""
}
