package http

import (
	"net/http"

	"github.com/lucat1/magik"
)

// Handles requests to the */reigster.. path authenticating
// the user who provides a valid token and redirecting the request appropriately
// TODO: set the cookie on the client and prob generate a different one with a different expiry date
func RegisterFunc(m *magik.Magik) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func Register(m *magik.Magik) http.Handler {
	return http.Handler(RegisterFunc(m))
}

// checks if the user is authenticated and sets the appropriate values in the context
// based on that. then it calls the h handler to properly manage the request
func WithAuthFunc(m *magik.Magik, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func WithAuth(m *magik.Magik, h http.Handler) http.Handler {
	return http.Handler(WithAuthFunc(m, h))
}
