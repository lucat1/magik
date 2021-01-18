package magik

import (
	"fmt"
	"net/url"
	"path"
	"time"
)

type MagikToken interface {
	Generate(email string, length time.Duration) (string, error)
	Validate(token string) (string, error)
}

type MagikEmail interface {
	// auth can be of varios types depending on the email provider
	// the Send method should check for the correctness of this type
	Send(auth interface{}, to, body string) error
}

type MagikEmailBody func(token, url string) string
type MagikConfig struct {
	BaseURL      string
	TokenTime    time.Duration
	RegisterBody MagikEmailBody
	LoginBody    MagikEmailBody
	EmailAuth    interface{}
}

type Magik struct {
	Config MagikConfig

	BaseURL *url.URL
	Token   MagikToken
	Email   MagikEmail
}

// TokenURL generates an url appending to the base url
// a path and providing the token value as a query parameter
// example:
// TokenURL("register", "abc") = "baseURL/register?t=abc"
func (m Magik) TokenURL(method, token string) string {
	old := m.BaseURL.Path

	m.BaseURL.Path = path.Join(m.BaseURL.Path, method)
	m.BaseURL.Query().Add("t", token)
	res := m.BaseURL.String()

	m.BaseURL.Query().Del("t")
	m.BaseURL.Path = old

	return res
}

// internal APIs helper to call *Body functions for
// generating email bodies
func (m Magik) Body(kind, token string) string {
	url := m.TokenURL(kind, token)
	switch kind {
	case "register":
		return m.Config.RegisterBody(token, url)
	case "login":
		return m.Config.LoginBody(token, url)
	}

	return "invalid body kind, this error should be reported to the administrator\n\n" + url
}

func NewMagik(conf MagikConfig, tok MagikToken, email MagikEmail) (*Magik, error) {
	var (
		u   *url.URL
		err error
	)

	if u, err = url.Parse(conf.BaseURL); conf.BaseURL == "" || err != nil {
		return nil, fmt.Errorf("invalid url: '%s'", conf.BaseURL)
	}

	return &Magik{
		BaseURL: u,
		Config:  conf,
		Token:   tok,
		Email:   email,
	}, nil
}
