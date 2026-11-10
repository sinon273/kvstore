package protocol

import (
	"strconv"
)

func EncodeOK() []byte {
	msg := "+OK\r\n"
	return []byte(msg)
}

func EncodeError(msg string) []byte {
	return []byte("-ERR " + msg + "\r\n")
}

func EncodeValue(v string) []byte {
	return []byte(v + "\r\n")
}

func EncodeInt(n int) []byte {
	strN := strconv.Itoa(n)
	return []byte(strN + "\r\n")
}
