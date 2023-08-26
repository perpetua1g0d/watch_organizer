package handler

import (
	"fmt"
	"log"
	"strings"
)

func ParseTabPath(tabPath string) (err error, tabs []string) {
	if len(tabPath) == 0 {
		tabPath = "root"
	}
	if tabs = strings.Split(tabPath, "/"); len(tabs) == 0 {
		return fmt.Errorf("Failed to split tabpath"), nil
	}

	tabs = append([]string{"root"}, tabs...)
	return
}

func (h *Handler) CreateTabPath(tabPath string) {
	err, tabs := ParseTabPath(tabPath)
	if err != nil {
		log.Fatalf("Failed to parse tabpath: %s", err.Error())
	}
	fmt.Println("tabs:", tabs)

	if err = h.service.Tab.CreateTabPath(tabs); err != nil {
		log.Fatalf("Failed to create tabpath: %s", err.Error())
	}

	return
}
