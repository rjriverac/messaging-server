package token

import "time"

type Maker interface {
	CreateToken(id int64, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
