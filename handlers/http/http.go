package http

import (
	"net/http"

	"github.com/lucat1/magik"
)

func RegisterFunc(m *magik.Magik) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("t")
		backto := r.URL.Query().Get("r")
		// email
		_, err := m.Token.Validate(token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: accept the email internally
		http.Redirect(w, r, backto, http.StatusMovedPermanently)
	})
}

func Register(m *magik.Magik) http.Handler {
	return http.Handler(RegisterFunc(m))
}
