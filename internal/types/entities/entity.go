package entities

import "database/sql"

type Entity struct {
	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	RawName sql.NullString `json:"raw_name"`
}
