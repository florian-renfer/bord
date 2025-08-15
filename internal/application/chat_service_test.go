package application

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/florian-renfer/b0red/internal/domain"
)

const (
	defaultUsername       = "Max Mustermann"
	secondaryUsername     = "Mia Musterfrau"
	defaultMessageContent = "This is a test message"
)

func TestCountConnections(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	chatService := NewChatService(logger)
	chatService.users = make(map[*domain.User]chan<- domain.Message)

	t.Run("Empty connections", func(t *testing.T) {
		count := chatService.CountConnections()
		assert.Equal(t, count, uint(0))
	})

	t.Run("One connection", func(t *testing.T) {
		chatService.users[&domain.User{}] = nil
		count := chatService.CountConnections()
		assert.Equal(t, count, uint(1))
	})

	t.Run("Multiple connections", func(t *testing.T) {
		chatService.users[&domain.User{}] = nil
		count := chatService.CountConnections()
		assert.Equal(t, count, uint(2))
	})
}

func TestRegisterConnection(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	chatService := NewChatService(logger)

	t.Run("Nil User", func(t *testing.T) {
		err := chatService.RegisterConnection(nil, nil)
		assert.Error(t, err)
	})

	t.Run("Nil outbound channel", func(t *testing.T) {
		user := &domain.User{
			Name: defaultUsername,
		}
		err := chatService.RegisterConnection(user, nil)
		assert.Error(t, err)
	})

	t.Run("Duplicate user", func(t *testing.T) {
		user := &domain.User{
			Name: defaultUsername,
		}
		chatService.RegisterConnection(user, make(chan<- domain.Message))
		err := chatService.RegisterConnection(user, make(chan<- domain.Message))
		assert.Error(t, err)
	})

	t.Run("One user", func(t *testing.T) {
		user := &domain.User{
			Name: defaultUsername,
		}
		err := chatService.RegisterConnection(user, make(chan<- domain.Message))
		assert.NoError(t, err)
	})
}

func TestUnregisterConnection(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	chatService := NewChatService(logger)

	t.Run("Nil user", func(t *testing.T) {
		err := chatService.UnregisterConnection(nil)
		assert.Error(t, err)
	})

	t.Run("Nil connections", func(t *testing.T) {
		err := chatService.UnregisterConnection(&domain.User{})
		assert.Error(t, err)
	})

	t.Run("Empty connections", func(t *testing.T) {
		err := chatService.UnregisterConnection(&domain.User{})
		assert.Error(t, err)
	})

	t.Run("One connection", func(t *testing.T) {
		u := &domain.User{
			Name: defaultUsername,
		}

		chatService.users = make(map[*domain.User]chan<- domain.Message)
		chatService.users[u] = make(chan<- domain.Message)

		assert.Equal(t, len(chatService.users), 1)
		err := chatService.UnregisterConnection(u)

		assert.NoError(t, err)
		assert.Equal(t, len(chatService.users), 0)
	})

	t.Run("One connection, nil channel", func(t *testing.T) {
		u := &domain.User{
			Name: defaultUsername,
		}

		chatService.users = make(map[*domain.User]chan<- domain.Message)
		chatService.users[u] = nil

		assert.Equal(t, len(chatService.users), 1)
		err := chatService.UnregisterConnection(u)

		assert.NoError(t, err)
		assert.Equal(t, len(chatService.users), 0)
	})

	t.Run("Two connections", func(t *testing.T) {
		u := &domain.User{
			Name: defaultUsername,
		}
		v := &domain.User{
			Name: secondaryUsername,
		}

		chatService.users = make(map[*domain.User]chan<- domain.Message)
		chatService.users[u] = make(chan<- domain.Message)
		chatService.users[v] = make(chan<- domain.Message)

		assert.Equal(t, len(chatService.users), 2)
		err := chatService.UnregisterConnection(u)

		assert.NoError(t, err)
		assert.Equal(t, len(chatService.users), 1)
		assert.Contains(t, chatService.users, v)
	})
}

func TestHanldeIncomingMessage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	chatService := NewChatService(logger)

	t.Run("No users", func(t *testing.T) {
		u := &domain.User{
			Name: defaultUsername,
		}

		msg := domain.Message{
			Content:   defaultMessageContent,
			Sender:    u,
			Timestamp: time.Now().UTC(),
		}

		chatService.users = nil

		err := chatService.HandleIncomingMessage(u, msg)
		assert.Error(t, err)
	})

	t.Run("Nil user", func(t *testing.T) {
		u := &domain.User{
			Name: defaultUsername,
		}

		msg := domain.Message{
			Content:   defaultMessageContent,
			Sender:    u,
			Timestamp: time.Now().UTC(),
		}

		chatService.users = make(map[*domain.User]chan<- domain.Message)

		err := chatService.HandleIncomingMessage(nil, msg)
		assert.Error(t, err)
	})

	t.Run("One user", func(t *testing.T) {
		u := &domain.User{
			Name: defaultUsername,
		}

		chatService.users = make(map[*domain.User]chan<- domain.Message)
		chatService.users[u] = make(chan<- domain.Message)

		msg := domain.Message{
			Content:   defaultMessageContent,
			Sender:    u,
			Timestamp: time.Now().UTC(),
		}

		err := chatService.HandleIncomingMessage(u, msg)
		assert.NoError(t, err)
	})

	t.Run("Two users", func(t *testing.T) {
		u := &domain.User{
			Name: defaultUsername,
		}

		v := &domain.User{
			Name: secondaryUsername,
		}

		chatService.users = make(map[*domain.User]chan<- domain.Message)
		chatService.users[u] = make(chan<- domain.Message, 1)
		chatService.users[v] = make(chan<- domain.Message, 1)

		msg := domain.Message{
			Content:   defaultMessageContent,
			Sender:    u,
			Timestamp: time.Now().UTC(),
		}

		err := chatService.HandleIncomingMessage(u, msg)
		assert.NoError(t, err)
	})
}
