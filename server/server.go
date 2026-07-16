package server

import (
	"net"

	"github.com/AnnaKhairetdinova/mini-redis/handler"
	"github.com/AnnaKhairetdinova/mini-redis/store"
)

func Start(port string, s *store.Store) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handler.Handle(conn, s)
	}

	return nil
}
