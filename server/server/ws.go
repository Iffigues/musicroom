package server

import (
)

type Ww struct {
	S *Server
}

func (ss *Server)NewWw(s *Server) (*Ww) {
	return &Ww{
		S: s,
	}
}

func (e *Ww) WWW(s *Server) {
}
