package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "subschema"
	app.Version = Version
	app.Usage = ""
	app.Author = "Akiomi Kamakura"
	app.Email = "akiomik@gmail.com"
	app.Action = doMain
	app.Run(os.Args)
}

func doMain(c *cli.Context) {
}
