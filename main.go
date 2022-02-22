package main

import (
	"github.com/xsadia/secred/cmd/app"
)

func main() {
	a := app.App{}

	a.Run(":1337")
}
