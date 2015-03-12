package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05-0700"
)

// NullTime replaces sql.NullTime with an implementation
// that supports proper JSON encoding/decoding.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// NewNullTime returns a new, properly instantiated
// NullTime object.
func NewNullTime(t time.Time) NullTime {
	return NullTime{Time: t, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *NullTime) Scan(value interface{}) error {
	ns.Time, ns.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Time, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullTime) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		if y := ns.Time.Year(); y < 0 || y >= 10000 {
			// RFC 3339 is clear that years are 4 digits exactly.
			// See golang.org/issue/4556#c15 for more discussion.
			return nil, errors.New("NullTime.MarshalJSON: year outside of range [0,9999]")
		}
		return []byte(ns.Time.Format(timeFormat)), nil
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullTime) UnmarshalJSON(text []byte) error {
	ns.Valid = false
	txt := string(text)
	if txt == "null" || txt == "" {
		return nil
	}

	// Fractional seconds are handled implicitly by Parse.
	t, err := time.Parse(timeFormat, string(txt))
	if err == nil {
		ns.Time = t
		ns.Valid = true
	}

	return err
}
