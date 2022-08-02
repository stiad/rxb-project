package api

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/rxbenefits/go-hw/moxbuster"
	"net/http"
	"strconv"
	"strings"
)

func (s *server) handleWelcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.LogDebug("handlers.Welcome")
		msg := struct {
			Message string `json:"message"`
		}{
			"Welcome to Moxbuster API",
		}
		s.respond(w, msg, http.StatusOK)
	}
}

func (s *server) handleListFilms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.LogDebug("handlers.handleListFilms")
		queryParams := s.getQueryParams(r)

		var title, rating, category string
		if len(queryParams["title"]) > 0 {
			title = "%" + queryParams["title"][0] + "%"
		}
		if len(queryParams["rating"]) > 0 {
			rating = strings.ToUpper(queryParams["rating"][0])
		}
		if len(queryParams["category"]) > 0 {
			category = "%" + strings.Title(queryParams["category"][0]) + "%"
		}

		rows, err := s.db.Query(
			"SELECT * FROM film_list WHERE ($1 = '' or title like $1) and ($2 = '' or rating = $2::mpaa_rating) and ($3 = '' or category like $3) limit 100",
			title, rating, category,
		)
		if err != nil {
			s.error(w, err, http.StatusInternalServerError)
			return
		}

		var films []moxbuster.Film
		for rows.Next() {
			film := moxbuster.Film{}
			err = rows.Scan(
				&film.FID,
				&film.Title,
				&film.Description,
				&film.Category,
				&film.Price,
				&film.Length,
				&film.Rating,
				&film.Actors,
			)
			films = append(films, film)
			if err != nil {
				s.error(w, err, http.StatusInternalServerError)
				return
			}
		}

		s.respond(w, films, http.StatusOK)
	}
}

func (s *server) handleGetFilm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.LogDebug("handlers.handleGetFilm")
		params := s.getParams(r)
		rows, err := s.db.Query("SELECT * FROM film_list WHERE fid = $1", params["filmId"])
		if err != nil {
			s.error(w, err, http.StatusInternalServerError)
			return
		}
		var films []moxbuster.Film
		for rows.Next() {
			film := moxbuster.Film{}
			err = rows.Scan(
				&film.FID,
				&film.Title,
				&film.Description,
				&film.Category,
				&film.Price,
				&film.Length,
				&film.Rating,
				&film.Actors,
			)
			films = append(films, film)
			if err != nil {
				s.error(w, err, http.StatusInternalServerError)
				return
			}
		}
		s.respond(w, films, http.StatusOK)
	}
}

func (s *server) handleListFilmComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.LogDebug("handlers.handleListFilmComments")
		params := s.getParams(r)
		rows, err := s.db.Query("SELECT * FROM film_comment WHERE fid = $1", params["filmId"])
		if err != nil {
			s.error(w, err, http.StatusInternalServerError)
			return
		}
		var comments []moxbuster.FilmComment
		for rows.Next() {
			comment := moxbuster.FilmComment{}
			err = rows.Scan(
				&comment.FirstName,
				&comment.LastName,
				&comment.CommentID,
				&comment.CommentBody,
				&comment.CreateDate,
				&comment.LastUpdate,
				&comment.FilmID,
				&comment.FilmTitle,
			)
			comments = append(comments, comment)
			if err != nil {
				s.error(w, err, http.StatusInternalServerError)
				return
			}
		}
		s.respond(w, comments, http.StatusOK)
	}
}

func (s *server) handleCreateNewComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := s.getParams(r)

		var comment struct {
			CustomerID int    `json:"customerID"`
			Body       string `json:"comment"`
		}

		err := s.decode(r, &comment)
		if err != nil {
			s.error(w, err, http.StatusBadRequest)
			return
		}

		comment.Body = (bluemonday.StripTagsPolicy()).Sanitize(comment.Body)

		filmId, _ := strconv.Atoi(params["filmId"])
		res, err := s.db.Exec("INSERT INTO comment (customer_id, body, film_id) VALUES ($1, $2, $3)", comment.CustomerID, comment.Body, filmId)
		if err != nil {
			s.error(w, err, http.StatusInternalServerError)
			return
		}

		if err != nil {
			s.error(w, err, http.StatusInternalServerError)
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			s.error(w, err, http.StatusInternalServerError)
			return
		}

		s.respond(w, id, http.StatusOK)
	}
}
