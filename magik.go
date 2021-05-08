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
  RegisterURL  string
	RegisterBody MagikEmailGenerator
  LoginURL     string
	LoginBody    MagikEmailGenerator
}

type Magik struct {
	Config MagikConfig

	BaseURL     *url.URL
	RegisterURL *url.URL
	LoginURL    *url.URL
	Token       MagikToken
	Email       MagikEmail
}

type EmailType uint
const (
  EmailTypeLogin    EmailType = iota
  EmailTypeRegister
)

// internal:
func absOrJoin(base, extra *url.URL) *url.URL {
  if(extra.IsAbs()) {
    return extra
  } else if(path.IsAbs(extra.Path)) {
    cpy := base
    cpy.Path = extra.Path
    return cpy
  } else {
    cpy := base
    cpy.Path = path.Join(base.Path, extra.Path)
    return cpy
  }
}

// TokenURL generates an url appending to the base url
// a path for the requested login type and providing the token value as
// a query parameter. Here's an example:
// TokenURL(EmailTypeRegister, "abc", "/home") = "baseURL/register?t=abc"
func (m Magik) TokenURL(kind EmailType, token, backto string) string {
  var extra *url.URL
  switch kind {
  case EmailTypeRegister:
    extra = m.RegisterURL
    break
  case EmailTypeLogin:
    extra = m.LoginURL
    break
  }
  url := absOrJoin(m.BaseURL, extra)
	q := url.Query()
	q.Set("t", token)
	q.Set("r", backto)
	url.RawQuery = q.Encode()

	return url.String()
}

// internal: helper to call *Body functions for
// generating email titles and bodies
func (m Magik) Body(kind EmailType, token, backto string) (string, string) {
	url := m.TokenURL(kind, token, backto)
	switch kind {
	case EmailTypeRegister:
		return m.Config.RegisterBody(token, url)
    break
	case EmailTypeLogin:
		return m.Config.LoginBody(token, url)
    break
	}

	return "attempt to authenticate", "invalid body kind, this error should be reported to the administrator\n\n" + url
}

func sendWithToken(m Magik, duration time.Duration, kind EmailType, email, backto string) error {
	token, err := m.Token.Generate(email, duration)
	if err != nil {
		return err
	}

	title, body := m.Body(kind, token, backto)
	return m.Email.Send(email, title, body)
}

func (m Magik) Register(email, backto string) error {
	return sendWithToken(m, m.Config.TokenTime, EmailTypeRegister, email, backto)
}

func (m Magik) Login(email, backto string) error {
	return sendWithToken(m, m.Config.TokenTime, EmailTypeLogin, email, backto)
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
		ur  *url.URL
		ul  *url.URL
		err error
	)

	if u, err = url.Parse(conf.BaseURL); conf.BaseURL == "" || err != nil {
		return nil, fmt.Errorf("invalid base url: '%s'", conf.BaseURL)
	}

	if ur, err = url.Parse(conf.RegisterURL); conf.RegisterURL == "" || err != nil {
		return nil, fmt.Errorf("invalid register url: '%s'", conf.RegisterURL)
	}

	if ul, err = url.Parse(conf.LoginURL); conf.LoginURL == "" || err != nil {
		return nil, fmt.Errorf("invalid login url: '%s'", conf.LoginURL)
	}

	return &Magik{
		BaseURL: u,
		RegisterURL: ur,
		LoginURL: ul,
		Config:  conf,
		Token:   tok,
		Email:   email,
	}, nil
}
