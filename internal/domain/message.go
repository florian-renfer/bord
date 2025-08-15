package domain

import "time"

type Message struct {
	Content   string
	Sender    *User
	Timestamp time.Time
}
