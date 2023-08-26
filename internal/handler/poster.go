package handler

import (
	"log"

	"github.com/perpetua1g0d/watch_organizer/internal/model"
)

func (h *Handler) CreatePoster(tabPath string, poster model.Poster) {
	err, tabs := ParseTabPath(tabPath)
	if err != nil {
		log.Fatalf("Failed to parse tabPath: %s", err.Error())
	}

	err = h.service.Tab.CreateTabPath(tabs)
	if err != nil {
		log.Fatalf("Failed to create tabpath: %s", err.Error())
	}

	err = h.service.Poster.Create(poster)
	if err != nil {
		log.Fatalf("Failed to create poster: %s", err.Error())
	}

	err, path := h.service.Tab.GetTabIds(tabs)
	if err != nil {
		log.Fatalf("Failed to get tab ids: %s", err.Error())
	}

	err = h.service.Tab.AddPosterToQueues(poster.Id, path)
	if err != nil {
		log.Fatalf("Failed to add poster to queues: %s", err.Error())
	}
}
