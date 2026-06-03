package errors

import "errors"

var (
	ErrUnsupportedType = errors.New("unsupported file type")
	ErrExternalTool    = errors.New("required external tool missing")
	ErrParseFailed     = errors.New("parse error")
)
