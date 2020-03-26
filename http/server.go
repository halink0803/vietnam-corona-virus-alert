package http

import (
	"github.com/gin-gonic/gin"
)

// Server type
type Server struct {
	r    *gin.Engine
	host string
}

// NewServer return new Server instance
func NewServer(host string) *Server {
	r := gin.Default()
	return &Server{
		r: r,
	}
}

// Run the server
func (s *Server) Run() error {
	return s.r.Run(s.host)
}
