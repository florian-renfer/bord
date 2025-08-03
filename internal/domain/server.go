package domain

type Server interface {
	ListenAndServe(addr string) error
	Broadcast(message Message) error
	Connect(user *User) error
	Disconnect(user *User) error
}
