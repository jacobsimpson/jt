package datetime

import (
	"fmt"
	"time"
)

type DateTimeFormat struct {
	Layout string
	// Used to indicate that the layout doesn't have a year in it, so assume
	// that the current year is the correct year. Sometimes dates are displayed
	// in a shortened format when the omitted information is the same as the
	// current information.
	UseCurrentYear   bool
	UseCurrentDay    bool
	UseLocalTimezone bool
}

// These are the formats that are valid literals when embedded unquoted in a jt
// script.
var LiteralFormats = []DateTimeFormat{
	{"2006-01-02T15:04:05.000Z", false, false, false},
	{"2006-01-02T15:04:05", false, false, true},
	{"2006-01-02T15:04", false, false, true},
	{"2006-01-02T15", false, false, true},
	{"2006-01-02T", false, false, true},
	{"01-02T", true, false, true},
	{"20060102T15:04:05.000Z", false, false, false},
	{"20060102T15:04:05", false, false, true},
	{"20060102T15:04", false, false, true},
	{"20060102T15", false, false, true},
	{"20060102T", false, false, true},
	{"0102T", true, false, true},
}

// These are the formats that will be attempted when an unknown value is
// compared to a date. The justification for the large degree of flexibilty is
// that by coomparing a value to a date type indicates the user is hoping for a
// valid date to be present.
var CoercionFormats = []DateTimeFormat{
	{"2006-01-02T15:04:05.000Z", false, false, false},
	{"2006-01-02T15:04:05", false, false, true},
	{"2006-01-02T15:04", false, false, true},
	{"2006-01-02T15", false, false, true},
	{"2006-01-02T", false, false, true},
	{"2006-01-02", false, false, true},
	{"20060102T15:04:05.000Z", false, false, false},
	{"20060102T15:04:05", false, false, true},
	{"20060102T15:04", false, false, true},
	{"20060102T15", false, false, true},
	{"20060102T", false, false, true},
	{"20060102", false, false, true},
	{"Mon Jan 2 15:04:05 PST 2006", false, false, false}, // Output of `date`
	{"Monday Jan 2 15:04:05 PST 2006", false, false, false},
	{"Mon Jan 2 15:04:05", true, false, true},
	{"Monday Jan 2 15:04:05", true, false, true},
	{"2Jan06", false, false, true}, // date/time in `ps -ef`
	{"15:05AM", true, true, true},  // date/time in `ps -ef`
	{"Jan 2 15:04", true, false, true},
	{"Jan 2, 2006", false, false, true},
	{"January 2, 2006", false, false, true},
	{"Jan _2 2006", false, false, true}, // Older files in `ls -l`
}

func ParseDateTime(formats []DateTimeFormat, str string) (*time.Time, error) {
	for _, f := range formats {
		if t, err := time.Parse(f.Layout, str); err == nil {
			if f.UseCurrentYear {
				t = time.Date(time.Now().Year(),
					t.Month(),
					t.Day(),
					t.Hour(),
					t.Minute(),
					t.Second(),
					t.Nanosecond(),
					t.Location())
			}
			if f.UseCurrentDay {
				t = time.Date(t.Year(),
					time.Now().Month(),
					time.Now().Day(),
					t.Hour(),
					t.Minute(),
					t.Second(),
					t.Nanosecond(),
					t.Location())
			}
			if f.UseLocalTimezone {
				t = time.Date(t.Year(),
					t.Month(),
					t.Day(),
					t.Hour(),
					t.Minute(),
					t.Second(),
					t.Nanosecond(),
					time.Now().Location())
			}
			return &t, nil
		}
	}
	return nil, fmt.Errorf("Unable to convert %q to a date", str)
}
