package errorx

import (
	"strings"

	"github.com/Martindeeepdark/go-common/errorx/internal"
)

// StatusError is an interface for error with status code, you can
// create an error through New or WrapByCode and convert it back to
// StatusError through FromStatusError to obtain information such as
// error status code.
type StatusError interface {
	error
	Code() int32
	Msg() string
	IsAffectStability() bool
	Extra() map[string]string
}

// Option is used to configure an StatusError.
type Option = internal.Option

// KV creates a parameter option for error message placeholder replacement
func KV(k, v string) Option {
	return internal.Param(k, v)
}

// KVf creates a parameter option with formatted value
func KVf(k, v string, a ...any) Option {
	return internal.Paramf(k, v, a...)
}

// Extra creates an extra key-value option
func Extra(k, v string) Option {
	return internal.Extra(k, v)
}

// New get an error predefined in the configuration file by statusCode
// with a stack trace at the point New is called.
func New(code int32, options ...Option) error {
	return internal.NewByCode(code, options...)
}

// WrapByCode returns an error annotating err with a stack trace
// at the point WrapByCode is called, and the status code.
func WrapByCode(err error, statusCode int32, options ...Option) error {
	if err == nil {
		return nil
	}

	return internal.WrapByCode(err, statusCode, options...)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return internal.Wrapf(err, format, args...)
}

// ErrorWithoutStack removes stack trace from error message
func ErrorWithoutStack(err error) string {
	if err == nil {
		return ""
	}
	errMsg := err.Error()
	index := strings.Index(errMsg, "stack=")
	if index != -1 {
		errMsg = errMsg[:index]
	}
	return errMsg
}

// Register registers an error code with message
func Register(code int32, msg string, opts ...internal.RegisterOption) {
	internal.Register(code, msg, opts...)
}

// SetDefaultErrorCode sets the default error code
func SetDefaultErrorCode(code int32) {
	internal.SetDefaultErrorCode(code)
}

// WithAffectStability creates a register option for affect stability
func WithAffectStability(affectStability bool) internal.RegisterOption {
	return internal.WithAffectStability(affectStability)
}
