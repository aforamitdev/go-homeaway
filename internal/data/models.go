package data

import (
	"database/sql"
	"errors"
)

type Models struct {
	Movie MovieMode
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

func NewModel(db *sql.DB) Models {
	return Models{
		Movie: MovieMode{DB: db},
	}
}
