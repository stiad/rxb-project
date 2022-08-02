package api

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/rxbenefits/go-hw/moxbuster"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestServer_handleListFilms(t *testing.T) {
	db, sqlMock, _ := sqlmock.New()

	sqlMock.
		ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"fid", "title", "description", "category", "price", "length", "rating", "actors"}).
			AddRow(1, "Academy Dinosaur", "A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies", "Documentary", 0.99, 86, "PG", "Penelope Guiness, Christian Gable, Lucille Tracy, Sandra Peck, Johnny Cage, Mena Temple, Warren Nolte, Oprah Kilmer, Rock Dukakis, Mary Keitel").
			AddRow(2, "Ace Goldfinger", "A Astounding Epistle of a Database Administrator And a Explorer who must Find a Car in Ancient China", "Horror", 4.99, 48, "G", "Bob Fawcett, Minnie Zellweger, Sean Guiness, Chris Depp"),
	)

	r, w := httptest.NewRequest(http.MethodGet, "/films?title=Academy", nil), httptest.NewRecorder()
	NewServer(db).handleListFilms()(w, r)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != `[{"fid":1,"title":"Academy Dinosaur","description":"A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies","category":"Documentary","price":0.99,"length":86,"rating":"PG","actors":"Penelope Guiness, Christian Gable, Lucille Tracy, Sandra Peck, Johnny Cage, Mena Temple, Warren Nolte, Oprah Kilmer, Rock Dukakis, Mary Keitel"},{"fid":2,"title":"Ace Goldfinger","description":"A Astounding Epistle of a Database Administrator And a Explorer who must Find a Car in Ancient China","category":"Horror","price":4.99,"length":48,"rating":"G","actors":"Bob Fawcett, Minnie Zellweger, Sean Guiness, Chris Depp"}]`+"\n" {
		t.Errorf("handleListFilms expected  %v", string(data))
	}
}

func TestServer_handleGetFilm(t *testing.T) {
	db, sqlMock, _ := sqlmock.New()

	sqlMock.
		ExpectQuery("SELECT").
		WillReturnRows(
			sqlmock.NewRows([]string{"fid", "title", "description", "category", "price", "length", "rating", "actors"}).
				AddRow(1, "Academy Dinosaur", "A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies", "Documentary", 0.99, 86, "PG", "Penelope Guiness, Christian Gable, Lucille Tracy, Sandra Peck, Johnny Cage, Mena Temple, Warren Nolte, Oprah Kilmer, Rock Dukakis, Mary Keitel"),
		)

	r, w := httptest.NewRequest(http.MethodGet, "/films/1", nil), httptest.NewRecorder()
	NewServer(db).handleGetFilm()(w, r)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("failed to unmarshal result %v", err)
	}

	var films []moxbuster.Film
	if err = json.Unmarshal(data, &films); err != nil {
		t.Error(err)
	}

	expected := `[{"fid":1,"title":"Academy Dinosaur","description":"A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies","category":"Documentary","price":0.99,"length":86,"rating":"PG","actors":"Penelope Guiness, Christian Gable, Lucille Tracy, Sandra Peck, Johnny Cage, Mena Temple, Warren Nolte, Oprah Kilmer, Rock Dukakis, Mary Keitel"}]` + "\n"
	if string(data) != expected {
		t.Errorf("expected %v got %v", expected, string(data))
	}
}

func TestServer_handleCreateNewComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("SQLMOCK ERROR: %v", err)
	}
	api := NewServer(db)
	w := httptest.NewRecorder()

	type body struct {
		CustomerID int    `json:"customerID"`
		Comment    string `json:"comment"`
	}
	type requestDetails struct {
		method string
		target string
		body   body
	}
	tt := []struct {
		request requestDetails
		writer  *httptest.ResponseRecorder
	}{
		{
			request: requestDetails{
				method: http.MethodPost,
				target: "/films/1/comments",
				body:   body{1, "this is a test comment"},
			},
			writer: w,
		},
		{
			request: requestDetails{
				method: http.MethodPost,
				target: "/films/2/comments",
				body:   body{2, "this is a test comment"},
			},
			writer: w,
		},
	}

	for _, tc := range tt {
		tc.writer.Flush()
		mock.ExpectExec("INSERT INTO comment").
			WithArgs(tc.request.body.CustomerID, tc.request.body.Comment, tc.request.body.CustomerID).
			WillReturnResult(sqlmock.NewResult(int64(tc.request.body.CustomerID), 0))

		b, _ := json.Marshal(tc.request.body)

		request := httptest.NewRequest(tc.request.method, tc.request.target, bytes.NewReader(b))
		request = mux.SetURLVars(request, map[string]string{"filmId": strconv.Itoa(tc.request.body.CustomerID)})

		api.handleCreateNewComment()(tc.writer, request)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("a sql expectation was not met: %v", err.Error())
		}
	}
}

func TestServer_handleListFilmComments(t *testing.T) {
	db, sqlMock, _ := sqlmock.New()

	datetime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")

	sqlMock.
		ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"first_name", "last_name", "comment_id", "body", "create_date", "last_updated", "fid", "title"}).
			AddRow("Sam", "Test", 1, "This is a test comment", datetime, datetime, 1, "Some Movie").
			AddRow("Mike", "Test", 2, "This is a test comment", datetime, datetime, 2, "Some Movie"),
	)

	r, w := httptest.NewRequest(http.MethodGet, "/films/1/comments", nil), httptest.NewRecorder()
	NewServer(db).handleListFilmComments()(w, r)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != `[{"firstName":"Sam","lastName":"Test","filmID":1,"filmTitle":"Some Movie","commentID":"1","commentBody":"This is a test comment","createDate":"0001-01-01T00:00:00Z","lastUpdate":"0001-01-01T00:00:00Z"},{"firstName":"Mike","lastName":"Test","filmID":2,"filmTitle":"Some Movie","commentID":"2","commentBody":"This is a test comment","createDate":"0001-01-01T00:00:00Z","lastUpdate":"0001-01-01T00:00:00Z"}]`+"\n" {
		t.Errorf("handleListFilms expected  %v", string(data))
	}
}
