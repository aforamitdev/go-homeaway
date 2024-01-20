package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
	query := `INSERT INTO movies(title,year,runtime,genres) VALUES ($1,$2,$3,$4) RETURNING id,created_at,version`
	fmt.Println(m.DB.Stats())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	fmt.Println(query)

	row := m.DB.QueryRowContext(ctx, query, args...)

	err := row.Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

	if err != nil {
		fmt.Println(err)
	}

	return err

}

func (m MovieMode) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id,created_at,title,year,runtime,genres,version FROM movies where id=$1`

	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)
	fmt.Println(err)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil

}

func (m MovieMode) Update(movie *Movie) error {
	query := `UPDATE movies SET title=$1,year=$2,runtime=$3,genres=$4,version=uuid_generate_v4() WHERE id=$5 RETURNING version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return err

}

func (m MovieMode) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cancel()

	result, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

func (m MovieMode) GetAll(title string, genres []string, filter Filter) ([]*Movie, error) {
	query := `SELECT id,created_at,title,year, runtime, genres,version FROM movies ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	movies := []*Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)

		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}
