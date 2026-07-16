package handler

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/AnnaKhairetdinova/mini-redis/parser"
	"github.com/AnnaKhairetdinova/mini-redis/store"
)

func Handle(conn net.Conn, s *store.Store) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		line := scanner.Text()
		cmd, err := parser.Parse(line)
		if err != nil {
			fmt.Fprintf(os.Stdout, "ERR: %s", err)
			return
		}

		result := execute(cmd, s)

		fmt.Fprintln(conn, result)
	}
}

func execute(cmd *parser.Command, s *store.Store) string {
	switch cmd.Name {
	case parser.SET:
		s.Set(cmd.Key, cmd.Value, cmd.TTL)
		return "OK"

	case parser.GET:
		res, ok := s.Get(cmd.Key)
		if !ok {
			return "nil"
		}
		return res

	case parser.DEL:
		ok := s.Del(cmd.Key)
		if !ok {
			return "nil"
		}
		return "OK"

	case parser.KEYS:
		keys := s.Keys()
		if len(keys) == 0 {
			return "пусто"
		}
		return strings.Join(keys, " ")

	case parser.PING:
		return "PONG"
	}

	return "ERR"
}
