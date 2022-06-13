package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorExpiredToken = errors.New("token is expired")
	ErrInvalidToken   = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"uuid"`
	User      int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Expires   time.Time `json:"expires"`
}

func CreatePayload(id int64, duration time.Duration) (*Payload, error) {
	ID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        ID,
		User:      id,
		CreatedAt: time.Now(),
		Expires:   time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.Expires) {
		return ErrorExpiredToken
	}
	return nil
}
