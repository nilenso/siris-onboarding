package server

func (s *Server) routes() {
	s.Router.Get("/ping", s.Ping)
}
