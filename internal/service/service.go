package service

import (
	"github.com/perpetua1g0d/watch_organizer/internal/model"
	"github.com/perpetua1g0d/watch_organizer/internal/repository"
)

type Poster interface {
	Create(poster model.Poster) (err error)
}

type Tab interface {
	CreateTabPath(tabs []string) error
	AddPosterToQueues(posterId int, tabs []string) (err error)
}

type Service struct {
	Poster
	Tab
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Poster: NewPosterService(repo.Poster),
		Tab:    NewTabService(repo.Tab),
	}
}
