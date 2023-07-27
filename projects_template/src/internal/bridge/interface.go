package bridge

import "time"

type Date interface {
	Today() time.Time
	Now() time.Time
}
