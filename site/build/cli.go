package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DexterLB/protopit/site/builder"
	"github.com/urfave/cli"
)

func build() {
	builder.BuildVariant("en", "content")
}

func serve(address string) {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("output/"))))
	log.Printf("web server at %s.", address)
	http.ListenAndServe(address, nil)
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "builds the site",
			Action: func(c *cli.Context) error {
				build()
				return nil
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "builds the site and starts a web server in its directory",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bind",
					Value: ":8080",
					Usage: "bind address",
				},
			},
			Action: func(c *cli.Context) error {
				serve(c.String("bind"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
