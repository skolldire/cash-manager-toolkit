package tcp

import (
	"context"
	"fmt"
	"github.com/skolldire/cash-manager-toolkit/pkg/client/log"
	"net"
	"sync"
)

type service struct {
	port string
	log  log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		port: d.Port,
		log:  d.Log,
	}
}

func (s service) Init(f ProcessingFunc) {
	var wg sync.WaitGroup
	defer wg.Done()
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.log.Error(context.Background(), err, map[string]interface{}{"port": s.port, "message": "Error starting server"})
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"port": s.port,
					"message": "Error closing listener"})
		}
	}(listener)
	s.log.Info(context.Background(), "Server listening",
		map[string]interface{}{"port": s.port})
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"port": s.port,
					"message": "Error accepting connection",
					"client":  conn.RemoteAddr().String()})
			continue
		}
		s.log.Info(context.Background(), "New connection accepted on port",
			map[string]interface{}{"port": s.port,
				"client": conn.RemoteAddr().String()})
		go s.handleConnection(conn, f)
	}
}

func (s service) handleConnection(conn net.Conn, f ProcessingFunc) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error closing connection"})
			return
		}
	}(conn)
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			s.log.Error(context.Background(), err,
				map[string]interface{}{"message": "Error reading from connection"})
		}
		message := string(buf[:n])
		response, err := f(message)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error responding:", err)
			return
		}
	}
}
