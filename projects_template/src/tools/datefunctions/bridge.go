package datefunctions

import (
	"time"

	"github.com/jinzhu/now"
)

type DateTool struct {
}

func NewDateTool() *DateTool {
	return &DateTool{}
}

func (u *DateTool) Today() time.Time {
	return now.BeginningOfDay()
}

func (u *DateTool) Now() time.Time {
	return time.Now()
}
