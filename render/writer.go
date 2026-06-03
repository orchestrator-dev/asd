package render

import (
	"asd/pager"
	"io"
	"reflect"
)

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

// NewWriter wraps the base writer with decorators based on Options.
func NewWriter(w io.Writer, opts any) io.WriteCloser {
	flat := false
	if reflect.ValueOf(opts).Kind() == reflect.Bool {
		flat = opts.(bool)
	} else {
		val := reflect.ValueOf(opts)
		if val.Kind() == reflect.Struct {
			flatField := val.FieldByName("Flat")
			if flatField.IsValid() {
				flat = flatField.Bool()
			}
		}
	}

	if !flat {
		w = pager.New(w)
		return w.(io.WriteCloser)
	}
	return nopCloser{w}
}
