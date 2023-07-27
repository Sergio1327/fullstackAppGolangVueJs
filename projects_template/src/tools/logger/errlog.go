package logger

import (
	"fmt"
	"path"
	"runtime"

	"github.com/fatih/color"
)

// ErrLog выводит ошибку с названием файла и номером строки
func ErrLog(err error) error {
	var (
		fileName string
		fileLine int
	)

	pc := make([]uintptr, 3)
	cnt := runtime.Callers(2, pc)
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		file, line := fu.FileLine(pc[i] - 1)
		fileName = fmt.Sprintf("%s/%s", path.Base(path.Dir(file)), path.Base(file))
		fileLine = line
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	return fmt.Errorf("%s \t %s=%d %s=%s", err.Error(), cyan("line"), fileLine, cyan("file"), fileName)
}

// ErrLogWithLogin выводит ошибку с названием файла и номером строки
func ErrLogWithLogin(err error, login string) error {
	var (
		fileName string
		fileLine int
	)

	pc := make([]uintptr, 3)
	cnt := runtime.Callers(2, pc)
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		file, line := fu.FileLine(pc[i] - 1)
		fileName = fmt.Sprintf("%s/%s", path.Base(path.Dir(file)), path.Base(file))
		fileLine = line
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	return fmt.Errorf("%s login=%s \t %s=%d %s=%s", err.Error(), login, cyan("line"), fileLine, cyan("file"), fileName)
}
