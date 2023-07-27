package datefunctions

import (
	"fmt"
	"time"
)

// InTimeSpan проверяет входит ли время в указанный промежуток
func InTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// DateDiff diff date
type DateDiff struct {
	BeginDate time.Time
	EndDate   time.Time
}

// NewDateDiff конструктор
func NewDateDiff(b, e time.Time) DateDiff {
	return DateDiff{
		BeginDate: b,
		EndDate:   e,
	}
}

// Zero добавляет ноль в начале
func (df DateDiff) Zero(v int) string {
	var r string
	if v < 10 {
		r = fmt.Sprintf("0%d", v)
	} else {
		r = fmt.Sprintf("%d", v)
	}

	return r
}

func (df DateDiff) String() string {
	y, m, d, h, min, s := DiffDates(df.BeginDate, df.EndDate)

	hours := df.Zero(h)
	mins := df.Zero(min)
	secs := df.Zero(s)

	if y > 0 {
		return fmt.Sprintf("%dг. %d мес. %d дн. %s:%s:%s", y, m, d, hours, mins, secs)
	} else if m > 0 {
		return fmt.Sprintf("%d мес. %d дн. %s:%s:%s", m, d, hours, mins, secs)
	}
	return fmt.Sprintf("%d дн. %s:%s:%s", d, hours, mins, secs)
}

// DiffDays diff days
func DiffDays(a, b time.Time) int {
	_, _, day, _, _, _ := DiffDates(a, b)
	return day
}

// DiffDates diff dates
func DiffDates(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// MustParse обязательный парсинг
func MustParse(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return t
}

var timeZone = time.Now().Location()

// SetDateTimeZone установить зону
func SetDateTimeZone(date time.Time) time.Time {
	return date.In(timeZone)
}

// ChangeDateTimeZone установить зону
func ChangeDateTimeZone(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), timeZone)
}

const (
	// Sunday дни недели
	Sunday time.Weekday = iota
	// Monday дни недели
	Monday
	// Tuesday дни недели
	Tuesday
	// Wednesday дни недели
	Wednesday
	// Thursday дни недели
	Thursday
	// Friday дни недели
	Friday
	// Saturday дни недели
	Saturday
)

const (
	// OneDayMinutes минут в неделю
	OneDayMinutes = 1440
)

const (
	// DateTimeMask маска с датой и временем
	DateTimeMask = "02.01.2006 15:04:05"
	// DateTimeWithoutSecondMask маска с датой и временем без секунд
	DateTimeWithoutSecondMask = "02.01.2006 15:04"
	// DateTimeMaskLocale маска с датой и временем с зоной
	DateTimeMaskLocale = "02.01.2006 15:04:05 -0700"
	// DateMask маска с датой
	DateMask = "02.01.2006"
)

// NonWorkHours нерабочие часы
var NonWorkHours = []int{18, 19, 20, 21, 22, 23, 0, 1, 2, 3, 4, 5, 6, 7, 8}

// AddMinutesWithoutWeekend без выходных
func AddMinutesWithoutWeekend(min int, t time.Time) time.Time {
	for {
		if t.Add(time.Minute*time.Duration(min)).Weekday() == Sunday || t.Add(time.Minute*time.Duration(min)).Weekday() == Saturday {
			min += 10
		} else {
			t = t.Add(time.Minute * time.Duration(min))
			break
		}
	}

	return t
}

// AddMinutesWithoutWeekendAndNights без выходных и ночей
func AddMinutesWithoutWeekendAndNights(min int, t time.Time) time.Time {
	for {
		if t.Add(time.Minute*time.Duration(min)).Weekday() == Saturday || t.Add(time.Minute*time.Duration(min)).Weekday() == Sunday || intInSlice(t.Add(time.Minute*time.Duration(min)).Hour(), NonWorkHours) {
			min += 10
		} else {
			t = t.Add(time.Minute * time.Duration(min))
			break
		}
	}

	return t
}

// AddMinutesWithoutSunday без воскресенья
func AddMinutesWithoutSunday(min int, t time.Time) time.Time {
	for {
		if t.Add(time.Minute*time.Duration(min)).Weekday() == Sunday {
			min += 10
		} else {
			t = t.Add(time.Minute * time.Duration(min))
			break
		}
	}

	return t
}

// AddMinutesWithoutSundayNight без ночи воскресенья
func AddMinutesWithoutSundayNight(min int, t time.Time) time.Time {
	for {
		if t.Add(time.Minute*time.Duration(min)).Weekday() == Sunday || intInSlice(t.Add(time.Minute*time.Duration(min)).Hour(), NonWorkHours) {
			min += 10
		} else {
			t = t.Add(time.Minute * time.Duration(min))
			break
		}
	}

	return t
}

// AddMinutesWithoutNight без ночей
func AddMinutesWithoutNight(min int, t time.Time) time.Time {
	for {
		if intInSlice(t.Add(time.Minute*time.Duration(min)).Hour(), NonWorkHours) {
			min += 10
		} else {
			t = t.Add(time.Minute * time.Duration(min))
			break
		}
	}

	return t
}

// intInSlice int in slice
func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// MonthCountSince поиск количества месяцев до текущего момента
func MonthCountSince(comparedDate time.Time) int {
	var (
		// текущее время
		now = time.Now()
		// trunc текущего времени до месяца
		truncedNow = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		// trunc сравниваемой даты до месяца
		truncedComparedDate = time.Date(comparedDate.Year(), comparedDate.Month(), 1, 0, 0, 0, 0, now.Location())
		// месяц с которого начнется отсчет
		startMonth = comparedDate.Month()
		// итерируемый месяц
		nextMonth time.Month
		// количество месяцев
		monthCount int
	)

	// крутим луп до тех пор, пока сравниваемая дата не будет больше текущей
	for truncedComparedDate.Before(truncedNow) {
		// добавляем один месяц
		truncedComparedDate = truncedComparedDate.AddDate(0, 1, 0)

		// получаем следующий месяц после итерации
		nextMonth = truncedComparedDate.Month()

		// если следующий месяц не равен стартовому, то пополняем количество месяцев
		if nextMonth != startMonth {
			monthCount++
		}

		// сохраняем месяц на котором остановились в предыдущей итерации
		startMonth = nextMonth
	}

	return monthCount
}
