// Copyright (c) 2016, Sandeep Raju Prabhakar <me@sandeepraju.in>
// All rights reserved.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sandeepraju/iswho/src/iswho"
	"github.com/urfave/cli"
)

func main() {
	// configure the app help text.
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[options] {{end}}query

VERSION:
   {{.Version}}

AUTHOR(S):
   {{range .Authors}}{{ . }}
   {{end}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
	app := cli.NewApp()
	app.Name = "iswho"
	app.Usage = "minimal client for the whois directory service"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Sandeep Raju",
			Email: "me@sandeepraju.in",
		},
	}
	app.Version = iswho.VERSION
	app.Action = func(c *cli.Context) error {
		// check for the object
		if len(c.Args()) != 1 {
			errorMessage := "Error: invalid usage! See " +
				"`iswho help` for more details."
			fmt.Println(errorMessage)
			return errors.New(errorMessage)
		}
		// object is given, process the whois request now.
		iw := iswho.NewIsWho(
			&iswho.IsWhoConfig{
				Host:    "in.whois-servers.net",
				Port:    43,
				Verbose: false,
			},
		)
		results, err := iw.Search("sandeepraju.in")
		if err != nil {
			// something went wrong while searching.
			fmt.Fprintf(os.Stderr, "Error: %s\n")
			return err
		}
		// NOTE: this is probably bad that the entire
		// result is dumped into results var all at once
		// and it might be better to pass back a stream
		// scanner, but it should be OK for now.
		fmt.Println(results)

		return nil
	}
	app.Run(os.Args)
}
