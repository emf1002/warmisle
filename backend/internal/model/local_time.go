package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LocalTime wraps time.Time to control SQLite serialization format.
// modernc.org/sqlite uses time.String() which produces "+0000 UTC" suffix
// that SQLite's strftime cannot parse. This type stores as "YYYY-MM-DD HH:MM:SS".
type LocalTime time.Time

func (t LocalTime) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	return tt.UTC().Format("2006-01-02 15:04:05"), nil
}

func (t *LocalTime) Scan(src interface{}) error {
	if src == nil {
		*t = LocalTime{}
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		*t = LocalTime(v)
	case string:
		tt, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			tt, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return fmt.Errorf("cannot parse time %q: %w", v, err)
			}
		}
		*t = LocalTime(tt)
	default:
		return fmt.Errorf("unsupported Scan type for LocalTime: %T", src)
	}
	return nil
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return []byte("null"), nil
	}
	return tt.MarshalJSON()
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = LocalTime{}
		return nil
	}
	tt := time.Time(*t)
	if err := tt.UnmarshalJSON(data); err != nil {
		return err
	}
	*t = LocalTime(tt)
	return nil
}

// GormDataType implements GORM's DataType interface
func (LocalTime) GormDataType() string {
	return "datetime"
}

// FromTime converts time.Time to LocalTime
func FromTime(t time.Time) LocalTime {
	return LocalTime(t)
}

// FromTimePtr converts *time.Time to *LocalTime
func FromTimePtr(t *time.Time) *LocalTime {
	if t == nil {
		return nil
	}
	v := LocalTime(*t)
	return &v
}

// ToTime converts LocalTime to time.Time
func (t LocalTime) ToTime() time.Time {
	return time.Time(t)
}

// ToTimePtr converts *LocalTime to *time.Time
func ToTimePtr(t *LocalTime) *time.Time {
	if t == nil {
		return nil
	}
	v := time.Time(*t)
	return &v
}
