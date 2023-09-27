package types

import (
	"encoding/json"
	"time"
)

// NullTime - кастомный тип для time, который может принимать значание null при сериализации.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func NewNullTime(time time.Time) *NullTime {
	return &NullTime{
		Time:  time,
		Valid: true,
	}
}

func (s NullTime) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.Time)
}

//TODO: UnmarshalJSON

