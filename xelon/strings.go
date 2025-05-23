package xelon

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// Stringify attempts to create a string representation of Xelon types.
func Stringify(message any) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

// stringifyValue was graciously cargoculted from the go-protobuf library.
func stringifyValue(w io.Writer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		_, _ = w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		_, _ = fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		stringifySlice(w, v)
		return
	case reflect.Struct:
		stringifyStruct(w, v)
	default:
		if v.CanInterface() {
			_, _ = fmt.Fprint(w, v.Interface())
		}
	}
}

func stringifySlice(w io.Writer, v reflect.Value) {
	_, _ = w.Write([]byte{'['})
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			_, _ = w.Write([]byte{' '})
		}

		stringifyValue(w, v.Index(i))
	}

	_, _ = w.Write([]byte{']'})
}

func stringifyStruct(w io.Writer, v reflect.Value) {
	if v.Type().Name() != "" {
		_, _ = w.Write([]byte(v.Type().String()))
	}

	// special case for time.Time values
	if v.Type() == timeType {
		_, _ = fmt.Fprintf(w, "{%s}", v.Interface())
		return
	}

	_, _ = w.Write([]byte{'{'})

	var sep bool
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			continue
		}
		if fv.Kind() == reflect.Slice && fv.IsNil() {
			continue
		}

		if sep {
			_, _ = w.Write([]byte(", "))
		} else {
			sep = true
		}

		_, _ = w.Write([]byte(v.Type().Field(i).Name))
		_, _ = w.Write([]byte{':'})
		stringifyValue(w, fv)
	}

	_, _ = w.Write([]byte{'}'})
}
