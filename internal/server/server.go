package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sinon273/kvstore/internal/protocol"
	"github.com/sinon273/kvstore/internal/storage"
)

func handle(conn net.Conn, store *storage.Store) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		cmd, err := protocol.Parse(line)
		if err != nil {
			fmt.Printf("Error parse str: %s\n", err)
			if _, err := conn.Write(protocol.EncodeError(err.Error())); err != nil {
				fmt.Printf("Error write str: %s\n", err)
				return
			}
			continue
		}
		switch cmd.Name {
		case protocol.CmdSet:
			store.Set(cmd.Args[0], cmd.Args[1], cmd.TTL)
			if _, err := conn.Write(protocol.EncodeOK()); err != nil {
				fmt.Printf("Error write str: %s\n", err)
				return
			}
		case protocol.CmdGet:
			val, ok := store.Get(cmd.Args[0])
			if ok {
				if _, err := conn.Write(protocol.EncodeValue(val)); err != nil {
					fmt.Printf("Error write str: %s\n", err)
					return
				}
			} else {
				if _, err := conn.Write(protocol.EncodeError("not found")); err != nil {
					fmt.Printf("Error write str: %s\n", err)
					return
				}
			}
		case protocol.CmdDelete:
			store.Delete(cmd.Args[0])
			if _, err := conn.Write(protocol.EncodeOK()); err != nil {
				fmt.Printf("Error write str: %s\n", err)
				return
			}
		case protocol.CmdTTL:
			ttl := store.TTL(cmd.Args[0])
			if _, err := conn.Write(protocol.EncodeInt(ttl)); err != nil {
				fmt.Printf("Error write str: %s\n", err)
				return
			}
		}

	}
}

func Run(ctx context.Context, addr string, store *storage.Store) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	var wg sync.WaitGroup

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	store.StartCleanup(ctx)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handle(conn, store)
		}()
	}

	wg.Wait()
	return nil
}
