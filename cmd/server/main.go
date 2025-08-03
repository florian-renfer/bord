package main

import (
	"log/slog"
	"os"

	"github.com/florian-renfer/b0red/internal/application"
	"github.com/florian-renfer/b0red/internal/infrastructure"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	userRegistry := application.UserRegistry{}

	server := infrastructure.NewTCPServer(logger, &userRegistry)
	application.NewChat(server)
}
