package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

//	type Movie struct {
//		ID        int64     `json:"id"`
//		CreatedAt time.Time `json:"-"`
//		Title     string    `json:"title,omitempty"`
//		Year      int32     `json:"year,omitempty"`
//		Runtime   Runtime   `json:"runtime,omitempty"`
//		Genres    []string  `json:"genres,omitempty"`
//		Version   int32     `json:"version"`
//	}
type MovieMode struct {
	DB *sql.DB
}

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"-"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func (m Movie) MarshalJSON() ([]byte, error) {

	var runtime string

	if m.Runtime != 0 {
		runtime = fmt.Sprintf("%d mins", m.Runtime)
	}

	type MovieAlias Movie

	aux := struct {
		MovieAlias
		Runtime string `json:"runtime,omitempty"`
	}{
		MovieAlias: MovieAlias(m),
		Runtime:    runtime,
	}
	return json.Marshal(aux)

}

func (m MovieMode) Insert(movie *Movie) error {

	query := `INSERT INTO movies(title,year,runtime,genres) VALUE ($1,$2,$3,$4) RETURN id,created_at,version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

}

func (m MovieMode) get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieMode) Update(movie *Movie) error {
	return nil
}

func (m MovieMode) Delete(id int64) error {
	return nil
}
