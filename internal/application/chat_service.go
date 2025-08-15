package application

import (
	"errors"
	"log/slog"

	"github.com/florian-renfer/b0red/internal/domain"
)

type Chat struct {
	logger *slog.Logger
	users  map[*domain.User]chan<- domain.Message
}

type ChatService interface {
	HandleIncomingMessage(user *domain.User, message domain.Message) error
	RegisterConnection(user *domain.User, outbound chan<- domain.Message) error
	UnregisterConnection(user *domain.User) error
	CountConnections() uint
}

func NewChatService(logger *slog.Logger) *Chat {
	return &Chat{
		logger: logger,
	}
}

func (c *Chat) HandleIncomingMessage(user *domain.User, message domain.Message) error {
	if c.users == nil {
		msg := "Could not handle incoming message. Users are not initialized."
		c.logger.Warn(msg, "users", c.users)
		return errors.New(msg)
	}

	if user == nil {
		msg := "Could not handle incoming message. Origin is nil."
		c.logger.Error(msg, "origin", user)
		return errors.New(msg)
	}

	for u, outbound := range c.users {
		if u != user {
			outbound <- message
		} else {
			c.logger.Debug("Broadcasting detected. Did not send the message to the outbound channel of origin.")
		}
	}

	return nil
}

func (c *Chat) RegisterConnection(user *domain.User, outbound chan<- domain.Message) error {
	if c.users == nil {
		c.logger.Debug("Users not present. Initializing.")
		c.users = make(map[*domain.User]chan<- domain.Message)
	}

	if user == nil || outbound == nil {
		msg := "Could not register user. Either the user or the outbound channel is nil."
		c.logger.Error(msg, "user", user, "outboundChannel", outbound)
		return errors.New(msg)
	}

	_, ok := c.users[user]
	if ok {
		msg := "Could not register user. A user with the given username is registered already."
		c.logger.Error(msg, "user", user.Name)
		return errors.New(msg)
	}

	c.users[user] = outbound
	c.logger.Debug("User registered", "user", user.Name)
	return nil
}

func (c *Chat) UnregisterConnection(user *domain.User) error {
	if user == nil {
		msg := "Can not unregister connection for user."
		c.logger.Error(msg, "user", user)
		return errors.New(msg)
	}

	if c.users == nil {
		msg := "Can not unregister connection for user, connections are not tracked."
		c.logger.Error(msg)
		return errors.New(msg)
	}

	outboundChannel, ok := c.users[user]
	if ok && outboundChannel != nil {
		close(outboundChannel)
	}

	delete(c.users, user)
	return nil
}

func (c *Chat) CountConnections() uint {
	return uint(len(c.users))
}
