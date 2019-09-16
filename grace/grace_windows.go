package grace

import (
	"log"
	"net"
	"net/http"
)

type Grace struct {
	Server   *http.Server
	Listener net.Listener
}

func New(addr string, handler http.Handler, graceful bool) (*Grace, error) {
	var ln net.Listener
	var err error
	ln, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := &http.Server{
		Handler: handler,
	}
	return &Grace{
		Server:   s,
		Listener: ln,
	}, nil
}

func (g *Grace) Start() {
	if err := g.Server.Serve(g.Listener); err != nil {
		log.Println(err)
	}
}
