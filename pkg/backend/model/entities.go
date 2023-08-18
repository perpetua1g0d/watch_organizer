package model

import "time"

type Poster struct {
	Id          int // `json:"id" db:"id"`
	KpLink      string
	Rating      float32
	Name        string
	Year        int
	CreatedDate time.Time
}

type PosterGenre struct {
	Id    int
	Genre string
}

type Tab struct {
	Id   int
	Name string
}

type TabChildren struct {
	Id1 int
	Id2 int
}

type TabQueue struct {
	TabId    int
	PosterId int
	Position int
}

type TabPoster struct {
	TabId    int
	PosterId int
}
