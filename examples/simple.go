package main

import (
	"net/http"
	"time"
  "log"

	"github.com/lucat1/magik"
	"github.com/lucat1/magik/generators/jwt"
	mhttp "github.com/lucat1/magik/handlers/http"
	mlog "github.com/lucat1/magik/senders/log"
)

var (
	loginBody    = magik.StandardFormat("register to example.com", "to login press <a href=\"%url%\">this link</a>")
	registerBody = magik.StandardFormat("login onto example.com", "to register press <a href=\"%url%\">this link</a>")
)

func main() {
	config := magik.MagikConfig{
		BaseURL:      "http://localhost:3000/auth",
		TokenTime:    time.Hour * 6,
		RegisterBody: registerBody,
    RegisterURL:  "register",
		LoginBody:    loginBody,
    LoginURL:     "login",
	}

	generator := jwt.NewGenerator("a very secret secret")
	// sender, err := smtp.NewSender(smtp.SMTPEmailConfig{
	// 	Email:    os.Getenv("EMAIL"),
	// 	Password: os.Getenv("PASSWORD"),
	// 	Hostname: "smtp.gmail.com",
	// 	Port:     587,
	// })
  sender, err := mlog.NewSender(log.Default())
	if err != nil {
		panic(err)
	}
	auth, err := magik.NewMagik(config, generator, sender)
	if err != nil {
		panic(err)
	}
	if err := auth.Register("user@domain.com", "/"); err != nil {
		panic(err)
	}

	http.Handle("/auth/register", mhttp.Register(auth))
	http.ListenAndServe(":3000", nil)
}
