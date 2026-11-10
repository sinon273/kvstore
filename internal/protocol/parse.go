package protocol

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Parse(line string) (Command, error) {
	line = strings.TrimRight(line, "\r\n")
	if strings.TrimSpace(line) == "" {
		return Command{}, fmt.Errorf("empty command")
	}

	tokens := strings.Fields(line)
	name := strings.ToUpper(tokens[0])
	args := tokens[1:]

	switch name {
	case CmdSet:
		if len(args) == 2 {
			return Command{Name: CmdSet, Args: args}, nil
		}

		if len(args) == 3 {
			n, err := strconv.Atoi(args[2])
			if err != nil {
				return Command{}, fmt.Errorf("invalid TTL: %s", err)
			}
			ttl := time.Duration(n) * time.Second
			return Command{
				Name: CmdSet,
				Args: args[:2],
				TTL:  ttl,
			}, nil
		}
		return Command{}, fmt.Errorf("SET requires 2 or 3 arguments")

	case CmdGet:
		if len(args) != 1 {
			return Command{}, fmt.Errorf("GET requires 1 argument, got %d", len(args))
		}
		return Command{Name: CmdGet, Args: args}, nil

	case CmdDelete:
		if len(args) != 1 {
			return Command{}, fmt.Errorf("DELETE requires 1 argument, got %d", len(args))
		}
		return Command{Name: CmdDelete, Args: args}, nil

	case CmdTTL:
		if len(args) != 1 {
			return Command{}, fmt.Errorf("TTL requires 1 argument, got %d", len(args))
		}
		return Command{
			Name: CmdTTL,
			Args: args,
		}, nil

	default:
		return Command{}, fmt.Errorf("unknown command: %s", name)
	}
}
