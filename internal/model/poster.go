package model

import "time"

type Poster struct {
	Id        int       `json:"id" db:"id"`
	KpLink    string    `db:"kplink"`
	Rating    float32   `db:"rating"`
	Name      string    `db:"name"`
	Year      int       `db:"year"`
	CreatedAt time.Time `db:"createdat"`
	Genres    []string
}

type PosterGenre struct {
	Id    int
	Genre string
}

type PosterInput struct {
	
}