package service

import "github.com/perpetua1g0d/watch_organizer/internal/repository"

type TabService struct {
	repo repository.Tab
}

func NewTabService(repo repository.Tab) *TabService {
	return &TabService{repo: repo}
}

func (s *TabService) CreateTabPath(tabs []string) error {
	return s.repo.CreateTabPath(tabs)
}
func (s *TabService) AddPosterToQueues(posterId int, tabs []string) (err error) {
	return s.repo.AddPosterToQueues(posterId, tabs)
}
