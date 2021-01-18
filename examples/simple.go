package main

import (
	"log"

	"github.com/lucat1/magik"
	"github.com/lucat1/magik/generators/jwt"
	"github.com/lucat1/magik/senders/stmp"
)

func main() {
	config := magik.MagikConfig{
		BaseURL: "/auth",
	}
	generator := jwt.NewGenerator("a very secret secret")
	email := stmp.NewEmail()
	auth, err := magik.NewMagik(config, generator, email)
	if err != nil {
		panic(err)
	}
	log.Println("a pointer", auth)
}
