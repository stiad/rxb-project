package moxbuster

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

type FilmRating string

type Film struct {
	FID         int        `json:"fid" validate:"required gte=1"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Price       float64    `json:"price"`
	Length      int        `json:"length"`
	Rating      FilmRating `json:"rating" validate:"required oneof=G PG PG-13 R NC-17"`
	Actors      string     `json:"actors"`
}

type FilmComment struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	FilmID      int       `json:"filmID"`
	FilmTitle   string    `json:"filmTitle"`
	CommentID   string    `json:"commentID"`
	CommentBody string    `json:"commentBody"`
	CreateDate  time.Time `json:"createDate"`
	LastUpdate  time.Time `json:"lastUpdate"`
}
