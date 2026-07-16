package parser

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type CommandName string

const (
	SET  CommandName = "SET"
	GET  CommandName = "GET"
	DEL  CommandName = "DEL"
	KEYS CommandName = "KEYS"
	PING CommandName = "PING"
	INFO CommandName = "INFO"
)

func (c CommandName) String() string {
	return string(c)
}

type Command struct {
	Name  CommandName
	Key   string
	Value string
	TTL   time.Duration
}

func Parse(line string) (*Command, error) {
	lineSplit := strings.Fields(line)
	if len(lineSplit) == 0 {
		return nil, errors.New("empty command")
	}

	name := CommandName(strings.ToUpper(lineSplit[0]))
	var key, value string
	var ttl time.Duration

	switch name {
	case SET:
		if len(lineSplit) < 3 {
			return nil, errors.New("SET requires key and value")
		}

		key = lineSplit[1]
		value = lineSplit[2]

		if len(lineSplit) >= 5 && lineSplit[3] == "EX" {
			i, err := strconv.Atoi(lineSplit[4])

			if err != nil || i <= 0 {
				return nil, err
			}

			ttl = time.Duration(i) * time.Second
		}
	case GET, DEL:
		if len(lineSplit) < 2 {
			return nil, errors.New(name.String() + " requires a key")
		}

		key = lineSplit[1]
	case KEYS, PING, INFO:
		// без аргументов
	default:
		return nil, errors.New("unknown command " + name.String())
	}

	com := &Command{
		name,
		key,
		value,
		ttl * time.Second,
	}

	return com, nil
}
