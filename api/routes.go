package api

func (s *server) routes() {
	s.registerRoute("/", s.handleWelcome(), []string{"GET"})
	s.registerRoute("/films", s.handleListFilms(), []string{"GET"})
	s.registerRoute("/films/{filmId}", s.handleGetFilm(), []string{"GET"})
	s.registerRoute("/films/{filmId}/comments", s.handleListFilmComments(), []string{"GET"})
	s.registerRoute("/films/{filmId}/comments", s.handleCreateNewComment(), []string{"POST"})
}
