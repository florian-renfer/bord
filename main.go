package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/florian-renfer/b0red/internal/domain"
	"github.com/florian-renfer/b0red/internal/usecase/usecase"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	messages := make(chan domain.BrokerEvent, 256)
	connections := make([]domain.Client, 0)

	broker := &domain.Broker{
		Connections: connections,
		Events:      messages,
		Logger:      logger,
	}

	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		logger.Error("Failed to start listener", "error", err)
		os.Exit(1)
	}
	defer listener.Close()

	go broker.Listen()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Failed to accept connection", "error", err)
			continue
		}

		logger.Info("Accepted connection", "remote_addr", conn.RemoteAddr())

		go usecase.BroadcastHandleConnection(conn, broker.Events, logger)
	}
}
