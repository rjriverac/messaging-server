package token

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < 32 {
		return nil, fmt.Errorf("invalid key size: key must be exactly %d characters in length", 32)
	}

	pKey, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		// paseto:       paseto.V4Local,
		symmetricKey: pKey,
	}
	return maker, nil
}
func (maker *PasetoMaker) CreateToken(id int64, duration time.Duration) (string, error) {
	payload, err := CreatePayload(id, duration)
	if err != nil {
		return "", err
	}

	marshalled, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// var unmarshalled *Payload
	// if err := json.Unmarshal(marshalled, &unmarshalled); err != nil {
	// 	fmt.Printf("error in unmarshaldebug\n")
	// }

	// fmt.Printf("marshalled payload: %+v\n", unmarshalled)
	token, err := paseto.NewTokenFromClaimsJSON(marshalled, nil)
	if err != nil {
		return "", err
	}

	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))
	token.SetNotBefore(time.Now())
	return token.V4Encrypt(maker.symmetricKey, nil), nil

}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParserForValidNow()

	payload, err := parser.ParseV4Local(maker.symmetricKey, token, nil)
	if err != nil {
		switch err.Error() {
		case "the ValidAt time is after this token expires":
			return nil, ErrorExpiredToken
		}
		return nil, ErrInvalidToken
	}

	var unmarshalled *Payload
	if err := json.Unmarshal(payload.ClaimsJSON(), &unmarshalled); err != nil {
		return nil, err
	}

	// fmt.Printf("payload in verify:%+v\n", unmarshalled)
	return unmarshalled, nil

}
