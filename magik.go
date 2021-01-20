package magik

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"
)

type MagikToken interface {
	Generate(email string, length time.Duration) (string, error)
	Validate(token string) (string, error)
}

type MagikEmail interface {
	// auth can be of varios types depending on the email provider
	// the Send method should check for the correctness of this type
	Send(to, title, body string) error
}

type MagikEmailGenerator func(token, url string) (string, string)
type MagikConfig struct {
	BaseURL      string
	TokenTime    time.Duration
	RegisterBody MagikEmailGenerator
	LoginBody    MagikEmailGenerator
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
func (m Magik) TokenURL(method, token, backto string) string {
	old := m.BaseURL.Path

	m.BaseURL.Path = path.Join(m.BaseURL.Path, method)
	q := m.BaseURL.Query()
	q.Set("t", token)
	q.Set("r", backto)
	m.BaseURL.RawQuery = q.Encode()
	res := m.BaseURL.String()

	m.BaseURL.RawQuery = ""
	m.BaseURL.Path = old

	return res
}

// internal APIs helper to call *Body functions for
// generating email bodies
func (m Magik) Body(kind, token, backto string) (string, string) {
	url := m.TokenURL(kind, token, backto)
	switch kind {
	case "register":
		return m.Config.RegisterBody(token, url)
	case "login":
		return m.Config.LoginBody(token, url)
	}

	return "attempt to authenticate", "invalid body kind, this error should be reported to the administrator\n\n" + url
}

func sendWithToken(m Magik, duration time.Duration, kind, email, backto string) error {
	token, err := m.Token.Generate(email, duration)
	if err != nil {
		return err
	}

	title, body := m.Body("register", token, backto)
	return m.Email.Send(email, title, body)
}

func (m Magik) Register(email, backto string) error {
	return sendWithToken(m, m.Config.TokenTime, "register", email, backto)
}

func (m Magik) Login(email, backto string) error {
	return sendWithToken(m, m.Config.TokenTime, "login", email, backto)
}

func StandardFormat(titleTemplate, bodyTemplate string) MagikEmailGenerator {
	return func(token, url string) (string, string) {
		return strings.ReplaceAll(strings.ReplaceAll(titleTemplate, "%token%", token), "%url%", url),
			strings.ReplaceAll(strings.ReplaceAll(bodyTemplate, "%token%", token), "%url%", url)
	}
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
