package types

import "time"

type MigrationObject struct {
	Key  string
	File string
}

type Migration struct {
	Id        int64
	Key       string
	IsApplied bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
