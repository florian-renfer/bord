package application

import (
	"github.com/florian-renfer/b0red/internal/domain"
)

type Chat struct {
	users map[*domain.User]chan<- domain.Message
}

type ChatService interface {
	HandleIncomingMessage(user *domain.User, message domain.Message)
	RegisterConnection(user *domain.User, outbound chan<- domain.Message)
	UnregisterConnection(user *domain.User)
}

func (c *Chat) HandleIncomingMessage(user *domain.User, message domain.Message) {
	for u, outbound := range c.users {
		if u != user {
			outbound <- message
		}
	}
}

func (c *Chat) RegisterConnection(user *domain.User, outbound chan<- domain.Message) {
	if c.users == nil {
		c.users = make(map[*domain.User]chan<- domain.Message)
	}
	c.users[user] = outbound
}

func (c *Chat) UnregisterConnection(user *domain.User) {
	close(c.users[user])
	delete(c.users, user)
}
