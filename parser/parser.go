package parser

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var CommandName string

const (
	SET  = "SET"
	GET  = "GET"
	DEL  = "DEL"
	KEYS = "KEYS"
	PING = "PING"
)

type Command struct {
	Name        string
	CommandName string
	Key         string
	Value       string
	TTL         time.Duration
}

func Parse(line string) (*Command, error) {
	lineSplit := strings.Fields(line)
	if len(lineSplit) == 0 {
		return nil, errors.New("empty command")
	}

	name := lineSplit[0]
	CommandName = name // что за глобальная переменная CommandName ?
	key := ""
	value := ""
	ttl := time.Duration(0)

	switch name {
	case SET:
		if len(lineSplit) < 3 {
			return nil, errors.New("SET requires key and value")
		}

		key = lineSplit[1]
		value = lineSplit[2]

		if len(lineSplit) >= 5 && lineSplit[3] == "EX" {
			i, err := strconv.Atoi(lineSplit[4])

			if err != nil {
				return nil, err
			}

			ttl = time.Duration(i)
		}
	case GET, DEL:
		if len(lineSplit) < 2 {
			return nil, errors.New(name + " requires a key")
		}

		key = lineSplit[1]
	case KEYS, PING:
		// без аргументов
	default:
		return nil, errors.New("unknown command" + name)
	}

	com := &Command{
		name,
		CommandName,
		key,
		value,
		ttl * time.Second,
	}

	return com, nil
}
