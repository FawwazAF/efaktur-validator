package http

func (s *Server) registerHandler() {
	s.router.GET("/", s.handler.Index.HandlerIndex)
	s.router.POST("/efaktur/validate", s.handler.Efaktur.HandlerValidateEfaktur)
}
