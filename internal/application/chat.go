package application

import (
	"github.com/florian-renfer/b0red/internal/domain"
)

type Chat struct {
	server domain.Server
}

func NewChat(server domain.Server) error {
	chat := Chat{
		server,
	}

	return chat.server.ListenAndServe("4000")
}
