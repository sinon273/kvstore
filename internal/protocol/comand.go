package protocol

import "time"

type Command struct {
	Name string
	Args []string
	TTL  time.Duration
}

const (
	CmdSet    = "SET"
	CmdGet    = "GET"
	CmdDelete = "DELETE"
	CmdTTL    = "TTL"
)
