package log

import (
	"projects_template/tools/sqlnull"
	"time"
)

type Row struct {
	Time          time.Time
	Flag          string
	Message       string
	SpecialFields map[string]SpecialField
	Module        string
	File          sqlnull.NullString
	Line          sqlnull.NullString
	Details       map[string]string
}

type SpecialField struct {
	Value string
	Type  string
}

func NewTestRow(details map[string]string) Row {
	return Row{
		Time:    time.Now(),
		Flag:    "DEBUG",
		Message: "test",
		SpecialFields: map[string]SpecialField{
			"c_id": {
				Value: "3378",
				Type:  "int",
			},
			"oper_login": {
				Value: "m.zarif@sarkor.uz",
				Type:  "string",
			},
			"se_id": {
				Value: "200189621",
				Type:  "int",
			},
		},
		Module:  "test",
		File:    sqlnull.NewString("test.go"),
		Line:    sqlnull.NewString("1234"),
		Details: details,
	}
}
