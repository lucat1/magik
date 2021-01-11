package main

import (
	"log"

	"github.com/lucat1/magik"
	"github.com/lucat1/magik/generators/jwt"
)

func main() {
	generator := jwt.NewGenerator("a very secret secret")
	auth := magik.NewMagik(generator)
	log.Println("a pointer", auth)
}
