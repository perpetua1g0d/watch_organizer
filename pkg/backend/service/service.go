package service

import "github.com/perpetua1g0d/watch_organizer/pkg/backend/repository"

type Service struct {
	
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}