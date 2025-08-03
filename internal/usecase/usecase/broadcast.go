package usecase

import (
	"bufio"
	"log/slog"
	"net"

	"github.com/florian-renfer/b0red/internal/domain"
)

func BroadcastHandleConnection(conn net.Conn, broker chan<- domain.BrokerEvent, logger *slog.Logger) {
	defer func() {
		logger.Info("Closing connection", "remote_addr", conn.RemoteAddr())
		broker <- domain.BrokerEvent{
			Type: domain.BrokerRemoveClient,
			Client: domain.Client{
				Conn: conn,
			},
		}
		conn.Close()
	}()

	conn.Write([]byte("Welcome to the chat server!\nWhat's your name > "))
	scanner := bufio.NewScanner(conn)

	var name string
	if scanner.Scan() {
		name = scanner.Text()
	} else {
		logger.Warn("No name provided, closing connection")
		return
	}

	client := domain.Client{
		Conn: conn,
		Name: name,
	}

	broker <- domain.BrokerEvent{
		Type:   domain.BrokerAddClient,
		Client: client,
	}

	for scanner.Scan() {
		text := scanner.Text()
		logger.Info("Received", "msg", text)
		event := domain.BrokerEvent{
			Client:  client,
			Message: text,
			Type:    domain.BrokerBroadcastMessage,
		}

		broker <- event
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Connection error", "error", err)
	}
}
