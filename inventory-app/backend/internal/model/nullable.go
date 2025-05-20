package model

import (
	"database/sql"
	"encoding/json"
)

type JSONNullFloat64 struct {
	sql.NullFloat64
}

func (n JSONNullFloat64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Float64)
}

func (n *JSONNullFloat64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		n.Float64 = 0
		return nil
	}
	var x float64
	if err := json.Unmarshal(b, &x); err != nil {
		n.Valid = false
		return err
	}
	n.Float64 = x
	n.Valid = true
	return nil
}
