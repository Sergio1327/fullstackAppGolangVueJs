package logger

import (
	"fmt"
	"path"
	"product_storage/internal/entity/log"
	"reflect"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogStorage interface {
	SaveLog(row log.Row) error
}

type DbHook struct {
	Module              string
	Storage             LogStorage
	SpecialFieldsFilter map[string]bool
}

func NewDBHook(storage LogStorage, module string, specialFields ...string) DbHook {
	specialFieldsFilter := make(map[string]bool)
	for _, row := range specialFields {
		specialFieldsFilter[row] = true
	}

	return DbHook{
		Module:              module,
		Storage:             storage,
		SpecialFieldsFilter: specialFieldsFilter,
	}
}

func (*DbHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}

var exceptionFields = map[string]bool{"file": true, "func": true, "line": true}

func Parse(entry *logrus.Entry, sf map[string]bool) (logDetails map[string]string, specialFields map[string]log.SpecialField) {

	logDetails = make(map[string]string)
	specialFields = make(map[string]log.SpecialField)

	for key, value := range entry.Data {
		k := strings.ToLower(key)

		if _, exists := sf[k]; !exists {
			if _, exists := exceptionFields[k]; !exists {
				logDetails[key] = fmt.Sprintf("%v", value)
			}
		} else {
			if _, exists := sf[k]; exists {
				specialFields[k] = log.SpecialField{
					Value: fmt.Sprintf("%v", value),
					Type:  fmt.Sprintf("%v", reflect.TypeOf(value)),
				}
			}
		}

	}

	return
}

func (d *DbHook) Fire(e *logrus.Entry) error {
	pc := make([]uintptr, 3)
	cnt := runtime.Callers(5, pc)
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			e.Data["file"] = path.Base(file)
			e.Data["line"] = line
			break
		}
	}

	module := strings.ToUpper(strings.Replace(strings.Replace(d.Module, ".log", "", -1), "_error", "", -1)) // имя модуля == название файла лога

	if module == "" {
		module = "NO MODULE"
	}

	var (
		logDetails    map[string]string
		specialFields map[string]log.SpecialField
	)

	if len(e.Data) > 0 {
		logDetails, specialFields = Parse(e, d.SpecialFieldsFilter)
	}

	// для json деталей
	// var (
	// err error
	// b   []byte
	// )

	// if len(logDetails) > 0 {
	// 	b, err = json.Marshal(&logDetails)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	row := log.Row{
		Time:          e.Time,
		Flag:          strings.ToUpper(e.Level.String()),
		Message:       e.Message,
		Module:        module,
		SpecialFields: specialFields,
		Details:       logDetails,
	}

	row.File.Scan(e.Data["file"])
	row.Line.Scan(e.Data["line"])

	return d.Storage.SaveLog(row)

}
