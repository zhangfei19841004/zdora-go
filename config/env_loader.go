package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	"github.com/fatih/structs"
)

func newLoader(prefix string) *loader {
	return &loader{CamelCase: true, DefaultTagName: "default", Prefix: prefix}
}

type loader struct {
	DefaultTagName string

	Prefix string

	CamelCase bool
}

func (l *loader) getPrefix(s *structs.Struct) string {
	if l.Prefix != "" {
		return l.Prefix
	}

	return s.Name()
}

func (l *loader) Load(s interface{}) error {
	strct := structs.New(s)
	prefix := l.getPrefix(strct)
	if l.DefaultTagName == "" {
		l.DefaultTagName = "default"
	}

	for _, field := range structs.Fields(s) {

		if err := l.processTagField(l.DefaultTagName, field); err != nil {
			return err
		}
		if err := l.processEnvField(prefix, field); err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) processTagField(tagName string, field *structs.Field) error {
	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := l.processTagField(tagName, f); err != nil {
				return err
			}
		}
	default:
		defaultVal := field.Tag(l.DefaultTagName)
		if defaultVal == "" {
			return nil
		}

		err := fieldSet(field, defaultVal)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *loader) processEnvField(prefix string, field *structs.Field) error {
	fieldName := l.getEnvFieldName(prefix, field.Name())

	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := l.processEnvField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		v := os.Getenv(fieldName)
		if v == "" {
			return nil
		}

		if err := fieldSet(field, v); err != nil {
			return err
		}
	}

	return nil
}

func (l *loader) PrintEnvs(s interface{}) {
	strct := structs.New(s)
	prefix := l.getPrefix(strct)
	for _, field := range structs.Fields(s) {
		l.printEnvField(prefix, field)
	}
}

func (l *loader) printEnvField(prefix string, field *structs.Field) {
	fieldName := l.getEnvFieldName(prefix, field.Name())

	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			l.printEnvField(fieldName, f)
		}
	default:
		fmt.Printf("  %s (default value: '%s')\n", fieldName, field.Tag(l.DefaultTagName))
	}
}

func (l *loader) getEnvFieldName(prefix string, name string) string {
	fieldName := strings.ToUpper(name)
	if l.CamelCase {
		fieldName = strings.ToUpper(strings.Join(camelcase.Split(name), "_"))
	}

	return strings.ToUpper(prefix) + "_" + fieldName
}

func fieldSet(field *structs.Field, v string) error {
	switch f := field.Value().(type) {
	case flag.Value:
		if v := reflect.ValueOf(field.Value()); v.IsNil() {
			typ := v.Type()
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}

			if err := field.Set(reflect.New(typ).Interface()); err != nil {
				return err
			}

			f = field.Value().(flag.Value)
		}

		return f.Set(v)
	}

	// TODO: add support for other types
	switch field.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}

		if err := field.Set(val); err != nil {
			return err
		}
	case reflect.Int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		if err := field.Set(i); err != nil {
			return err
		}
	case reflect.String:
		if err := field.Set(v); err != nil {
			return err
		}
	case reflect.Slice:
		switch t := field.Value().(type) {
		case []string:
			if err := field.Set(strings.Split(v, ",")); err != nil {
				return err
			}
		case []int:
			var list []int
			for _, in := range strings.Split(v, ",") {
				i, err := strconv.Atoi(in)
				if err != nil {
					return err
				}

				list = append(list, i)
			}

			if err := field.Set(list); err != nil {
				return err
			}
		default:
			return fmt.Errorf("multiconfig: field '%s' of type slice is unsupported: %s (%T)",
				field.Name(), field.Kind(), t)
		}
	case reflect.Float64:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		if err := field.Set(f); err != nil {
			return err
		}
	case reflect.Int64:
		switch t := field.Value().(type) {
		case time.Duration:
			d, err := time.ParseDuration(v)
			if err != nil {
				return err
			}

			if err := field.Set(d); err != nil {
				return err
			}
		case int64:
			p, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				return err
			}

			if err := field.Set(p); err != nil {
				return err
			}
		default:
			return fmt.Errorf("multiconfig: field '%s' of type int64 is unsupported: %s (%T)",
				field.Name(), field.Kind(), t)
		}

	default:
		return fmt.Errorf("multiconfig: field '%s' has unsupported type: %s", field.Name(), field.Kind())
	}

	return nil
}
