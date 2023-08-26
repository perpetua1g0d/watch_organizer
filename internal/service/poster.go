package service

import (
	"github.com/perpetua1g0d/watch_organizer/internal/model"
	"github.com/perpetua1g0d/watch_organizer/internal/repository"
)

type PosterService struct {
	repo repository.Poster
}

func NewPosterService(repo repository.Poster) *PosterService {
	return &PosterService{repo: repo}
}

func (s *PosterService) Create(poster model.Poster) (err error) {
	return s.repo.Create(poster)
}
