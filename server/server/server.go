package server

import (
	"net/http"
	"os"
	"polaroid/config"
	"polaroid/pk"
	"strings"

	"polaroid/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type HH interface {
	WWW(*Server)
}

type Server struct {
	Router *mux.Router
	Data   *types.Data
	Handle map[string]*Handle
	Give   []HH
}

type Handle struct {
	Role   int
	Route  string
	Method []string
	Handle http.Handler
}

func (s *Server) AddHH(p ...HH) {
	for _, val := range p {
		s.Give = append(s.Give, val)
	}
}

func (s *Server) StartHH() {
	for _, val := range s.Give {
		val.WWW(s)
	}
}
