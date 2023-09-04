package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/perpetua1g0d/watch_organizer/internal/model"
)

type Poster interface {
	Create(poster model.Poster) (err error)
}

type Tab interface {
	CreateTabPath(tabs []string) (err error)
	AddPosterToQueues(posterId int, tabs []string) (err error)
}

type Repository struct {
	Poster
	Tab
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Poster: NewPosterPostgres(db),
		Tab:    NewTabPostgres(db),
	}
}
