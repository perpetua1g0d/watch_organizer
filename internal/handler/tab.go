package handler

import (
	"log"
	"strings"
)

func ParseTabPath(tabPath string) (tabs []string) {
	tabs = strings.Split(tabPath, "/")
	tabs = append([]string{"root"}, tabs...)
	return
}

func (h *Handler) CreateTabPath(tabPath string) {
	tabs := ParseTabPath(tabPath)
	// fmt.Println("parsed tabs:", tabs)

	if err := h.service.Tab.CreateTabPath(tabs); err != nil {
		log.Fatalf("Failed to create tabpath: %s", err.Error())
	}

	return
}
