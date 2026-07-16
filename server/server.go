package server

import (
	"context"
	"net"
	"sync"

	"github.com/AnnaKhairetdinova/mini-redis/handler"
	"github.com/AnnaKhairetdinova/mini-redis/store"
)

func Start(ctx context.Context, port string, s *store.Store) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	var wg sync.WaitGroup

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			handler.Handle(conn, s)
		}()
	}
}
