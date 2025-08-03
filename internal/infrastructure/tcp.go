package infrastructure

import (
	"log/slog"
	"net"

	"github.com/florian-renfer/b0red/internal/application"
	"github.com/florian-renfer/b0red/internal/domain"
)

type TCPServer struct {
	logger       *slog.Logger
	userRegistry *application.UserRegistry
	connections  map[*domain.User]net.Conn
}

func NewTCPServer(logger *slog.Logger, userRegistry *application.UserRegistry) *TCPServer {
	return &TCPServer{
		logger:       logger,
		userRegistry: userRegistry,
		connections:  make(map[*domain.User]net.Conn),
	}
}

func (s *TCPServer) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("Failed to start TCP server", "error", err, "address", addr)
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Failed to accept connection", "error", err)
			continue
		}

		s.logger.Debug("Accepted connection", "remote_addr", conn.RemoteAddr())

		user := &domain.User{
			Name: conn.RemoteAddr().String(),
		}
		s.connections[user] = conn
		s.userRegistry.AddUser(user)
	}
}

func (s *TCPServer) Broadcast(message domain.Message) error {
	for _, user := range s.userRegistry.GetUsers() {
		connection, ok := s.connections[user]
		if ok {
			connection.Write([]byte(message.Conetent + "\n"))
		}

	}
	return nil
}

func (s *TCPServer) Connect(user *domain.User) error {
	return nil
}

func (s *TCPServer) Disconnect(user *domain.User) error {
	return nil
}
