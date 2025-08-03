package domain

import "log/slog"

type BrokerEventType int

const (
	BrokerAddClient BrokerEventType = iota
	BrokerRemoveClient
	BrokerBroadcastMessage
)

type BrokerEvent struct {
	Type    BrokerEventType
	Client  Client
	Message string
}

type Broker struct {
	Connections []Client
	Events      chan BrokerEvent
	Logger      *slog.Logger
}

func (b *Broker) Listen() {
	for {
		event := <-b.Events

		switch event.Type {
		case BrokerAddClient:
			b.Logger.Info("Adding connection", "remote_addr", event.Client.Conn.RemoteAddr())
			b.Connections = append(b.Connections, event.Client)
		case BrokerRemoveClient:
			b.Logger.Info("Removing connection", "remote_addr", event.Client.Conn.RemoteAddr())
			for i, conn := range b.Connections {
				if conn == event.Client {
					b.Connections = append(b.Connections[:i], b.Connections[i+1:]...)
					break
				}
			}
		case BrokerBroadcastMessage:
			b.Logger.Info("Broadcasting message", "msg", event.Message, "receivers", len(b.Connections))

			for _, client := range b.Connections {
				if client == event.Client {
					continue
				}

				if _, err := client.Conn.Write([]byte(event.Client.Name + " > " + event.Message + "\n")); err != nil {
					b.Logger.Error("Failed to send message", "error", err, "remote_addr", client.Conn.RemoteAddr())
				}
			}
		}
	}
}
