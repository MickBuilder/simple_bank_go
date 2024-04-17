package auth

import "time"

type Token interface {
	Create(username string, duration time.Duration) (string, error)
	Verify(token string) (*Payload, error)
}
