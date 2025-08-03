package domain

import "time"

type Message struct {
	Conetent  string
	Sender    *User
	Timestamp time.Time
}
