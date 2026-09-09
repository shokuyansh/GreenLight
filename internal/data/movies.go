package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.db.QueryRowContext(ctx, stmt, args...).Scan(&movie.ID, &movie.CreatedAT, &movie.Version)
}

func (m MovieModel) Get(id int) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `select id,created_at,title,year,runtime,genres,version
	from movies
	where id=$1`

	var movie Movie
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
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
	where id=$5 and version=$6
	returning version`

	args := []any{
		movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConfict
		default:
			return err
		}
	}
	return nil
}

func (m MovieModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `delete from movies
	where id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.db.ExecContext(ctx, query, id)
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

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	query := fmt.Sprintf(`select count(*) OVER(),id,created_at,title,year,runtime,genres,version
	from movies 
	where (to_tsvector('simple',title)@@plainto_tsquery('simple',$1) OR $1='')
	AND (genres@>$2 or $2='{}')
	order by %s %s,id ASC
	limit $3 offset $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAT,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return movies, metadata, nil
}
