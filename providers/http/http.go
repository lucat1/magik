package http

import (
	"github.com/lucat1/magik"
	"net/http"
)

func Register(m magik.Magik) http.Handler {
	return http.Handler(Register(m))
}

func RegisterFunc(m magik.Magik) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := "test@example.com"
		token, err := m.Token.Generate(email, m.Config.TokenTime)
		if err != nil {
			// TODO: handle error while creating a token, prob a 500 error
		}

		m.Email.Send(m.Config.EmailAuth, email, m.Body("register", token))
	}
}
