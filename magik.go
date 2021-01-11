package magik

import "time"

type MagikToken interface {
	Generate(email string, length time.Duration) (string, error)
	Validate(token string) (string, error)
}

type MagikAuth struct {
	token MagikToken
}

func NewMagik(tok MagikToken) MagikAuth {
	return MagikAuth{
		token: tok,
	}
}
