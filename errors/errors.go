package errors

import (
	"runtime"

	base_errors "github.com/asyou-me/lib.v1/errors"
)

// New 创建一个错误容器，不记录错误路径
func New(code int, values ...string) *base_errors.ErrStruct {
	return codes.New(code, values...)
}

// NewWithPath 创建一个错误容器，记录错误路径
func NewWithPath(code int, values ...string) *base_errors.ErrStruct {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	err := codes.New(code, values...)
	err.Log = &base_errors.LogStruct{
		Local: file + ":" + f.Name(),
		Line:  line,
	}
	return err
}
