package store

import (
	"bufio"
	"net"
	"strings"
)

func (s *Store) handleTCP(conn net.Conn) {
	defer conn.Close()
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		s.logger.Println("error reading:", err.Error())
		return
	}
	tmp := strings.TrimSpace(string(input))
	s.logger.Printf("new tcp msg: %s\n", tmp)

	// trim spaces
	cmd := strings.Fields(tmp)
	// handle command
	rsp := s.handleCmd(cmd)
	// send a response back
	conn.Write([]byte(rsp))
}

// Select the handler.
func (s *Store) handleCmd(cmd []string) string {
	if len(cmd) == 0 {
		return "ERROR"
	}
	verb := strings.ToLower(cmd[0])
	args := cmd[1:]

	s.logger.Printf("processing %s command", cmd)
	switch verb {
	case "ping":
		return "PONG"
	case "join":
		return s.handleJoin(args)
	case "leave":
		return s.handleLeave(args)
	case "get":
		return s.handleGet(args)
	case "set":
		return s.handleSet(args)
	case "del":
		return s.handleDel(args)
	default:
		return "ERROR"
	}
}

func (s *Store) handleJoin(args []string) string {
	if len(args) != 2 {
		return "ERROR"
	}

	raftAddr := args[0]
	nodeID := args[1]
	if err := s.Join(nodeID, raftAddr); err != nil {
		return err.Error()
	}
	return "SUCCESS"
}

func (s *Store) handleLeave(args []string) string {
	if len(args) != 1 {
		return "ERROR"
	}

	nodeID := args[0]
	if err := s.Leave(nodeID); err != nil {
		return err.Error()
	}
	return "SUCCESS"
}

func (s *Store) handleGet(args []string) string {
	if len(args) != 1 {
		return "ERROR"
	}

	k := args[0]
	v, ok := s.Get(k)
	if !ok {
		return "ERROR"
	}
	return v
}

func (s *Store) handleSet(args []string) string {
	if len(args) != 2 {
		return "ERROR"
	}

	k := args[0]
	v := args[1]
	if err := s.Set(k, v); err != nil {
		return err.Error()
	}
	return "SUCCESS"
}

func (s *Store) handleDel(args []string) string {
	if len(args) != 1 {
		return "ERROR"
	}

	k := args[0]
	if err := s.Delete(k); err != nil {
		return err.Error()
	}
	return "SUCCESS"
}
