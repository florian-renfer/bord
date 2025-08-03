package infrastructure

import (
	"bufio"
	"log/slog"
	"net"
	"time"

	"github.com/florian-renfer/b0red/internal/application"
	"github.com/florian-renfer/b0red/internal/domain"
)

type TCPServer struct {
	logger      *slog.Logger
	chatService application.ChatService
	connections map[*domain.User]net.Conn
}

func NewTCPServer(logger *slog.Logger, chatService application.ChatService) *TCPServer {
	return &TCPServer{
		logger:      logger,
		chatService: chatService,
		connections: make(map[*domain.User]net.Conn),
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
		s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	user := &domain.User{}

	// FIXME: This leads to a closed connection because reading and writing is done in separate goroutines.
	// defer func() {
	// 	conn.Close()
	// 	s.chatService.UnregisterConnection(user)
	// }()

	// FIXME: This is a blocking call and should be handled in a separate goroutine.
	conn.Write([]byte("Welcome to the chat server!\nWhat's your name > "))
	scanner := bufio.NewScanner(conn)

	if scanner.Scan() {
		user.Name = scanner.Text()
		conn.Write([]byte(user.Name + " > "))
	} else {
		s.logger.Warn("No name provided, closing connection")
		return
	}

	outbound := make(chan domain.Message, 16)
	s.chatService.RegisterConnection(user, outbound)

	go s.handeOutgoingMessage(outbound, conn)
	go s.handleIncomingMessage(scanner, user)
}

func (s *TCPServer) handleIncomingMessage(scanner *bufio.Scanner, user *domain.User) {
	for scanner.Scan() {
		text := scanner.Text()
		s.logger.Info("Received", "msg", text)

		message := domain.Message{
			Conetent:  text,
			Sender:    user,
			Timestamp: time.Now().UTC(),
		}

		s.chatService.HandleIncomingMessage(user, message)
	}
	if err := scanner.Err(); err != nil {
		s.logger.Error("Connection error", "error", err)
	}
}

func (s *TCPServer) handeOutgoingMessage(outbound <-chan domain.Message, conn net.Conn) {
	for message := range outbound {
		conn.Write([]byte(message.Sender.Name + ": " + message.Conetent + "\n"))
	}
}
