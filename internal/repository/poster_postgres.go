package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/perpetua1g0d/watch_organizer/internal/model"
)

type PosterPostgres struct {
	db *sqlx.DB
}

func NewPosterPostgres(db *sqlx.DB) *PosterPostgres {
	return &PosterPostgres{db: db}
}

func createGenresInDB(db *sqlx.DB, poster model.Poster) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin genres tx: %w", err)
	}

	for _, genre := range poster.Genres {
		err = db.QueryRowx("insert into postergenre values (?, ?);", poster.Id, genre).Err()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to insert genre %s: %w", genre, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit genres tx: %w", err)
	}
	return
}

func (r *PosterPostgres) Create(poster model.Poster) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin poster tx: %w", err)
	}

	var posterId int
	err = r.db.QueryRowx(fmt.Sprintf("insert into poster values(default, '%s', %f, '%s', %d) returning id;",
		poster.KpLink, poster.Rating, poster.Name, poster.Year)).Scan(&posterId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to insert poster %+v: %w", poster, err)
	}

	if err = createGenresInDB(r.db, poster); err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create genres: %w", err)
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("Failed to commit poster tx: %w", err)
	}
	return
}
