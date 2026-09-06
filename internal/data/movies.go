package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/shokuyansh/GreenLight/internal/validator"
)

type Movie struct {
	ID        int       `json:"id"`
	CreatedAT time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int       `json:"version"`
}

type MovieModel struct {
	db *sql.DB
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(len(movie.Title) != 0, "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int(time.Now().Year()), "year", "must not be in future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain atleast 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

func (m MovieModel) Insert(movie *Movie) error {
	stmt := `insert into movies(title,year,runtime,genres) 
	values($1,$2,$3,$4)
	returning id,created_at,version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return m.db.QueryRow(stmt, args...).Scan(&movie.ID, &movie.CreatedAT, &movie.Version)
}

func (m MovieModel) Get(id int) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `select id,created_at,title,year,runtime,genres,version
	from movies
	where id=$1`

	var movie Movie
	err := m.db.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAT,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `update movies
	set title=$1,year=$2,runtime=$3,genres=$4,version=version+1
	where id=$5
	returning version`

	args := []any{
		movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID,
	}

	return m.db.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `delete from movies
	where id=$1`
	result, err := m.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}
	if rowAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
