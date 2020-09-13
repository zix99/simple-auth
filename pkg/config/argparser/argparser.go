package argparser

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func setFromString(v *reflect.Value, val string) error {
	switch v.Kind() {
	case reflect.Int:
		ival, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(ival))
		return nil
	case reflect.Bool:
		ival, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(ival))
		return nil
	case reflect.String:
		v.Set(reflect.ValueOf(val))
		return nil
	case reflect.Slice:
		av := reflect.Append(*v, reflect.ValueOf(val))
		v.Set(av)
		return nil
	}
	return fmt.Errorf("Unparseable type of kind %v", v.Kind())
}

// setDeeply sets a struct object to a value given a named-path
func setDeeply(obj interface{}, val string, path ...string) error {
	v := reflect.Indirect(reflect.ValueOf(obj))
	for _, s := range path {
		v = v.FieldByNameFunc(func(fName string) bool {
			return strings.ToLower(fName) == s
		})
		if !v.CanSet() {
			return errors.New("Not settable")
		}
	}
	return setFromString(&v, val)
}

func parseFlag(obj interface{}, flag, val string) error {
	if len(flag) < 2 {
		return fmt.Errorf("Invalid flag: %s", flag)
	}
	if flag[:2] == "--" {
		parts := strings.Split(flag[2:], "-")
		return setDeeply(obj, val, parts...)
	} else if flag[0] == '-' {
		parts := strings.Split(flag[1:], "-")
		return setDeeply(obj, "true", parts...)
	} else {
		return fmt.Errorf("Invalid flag: %s", flag)
	}
}

// LoadArgs nested obj (struct) given arguments in format --key=val
func LoadArgs(obj interface{}, args ...string) error {
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		key := parts[0]
		val := "true"
		if len(parts) > 1 {
			val = strings.Join(parts[1:], "=")
		}
		if err := parseFlag(obj, key, val); err != nil {
			return err
		}
	}
	return nil
}
