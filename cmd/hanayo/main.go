package main

import (
	"fmt"
	nhttp "net/http"
	"os"

	"git.zxq.co/ripple/hanayo"
	"git.zxq.co/ripple/hanayo/helpers/conf"
	"git.zxq.co/ripple/hanayo/http"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	hc := os.Getenv("HANAYO_CONFIG")
	if err := conf.Load(hc); err != nil {
		fmt.Println("Writing config")
		err = conf.Export(hc)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	app := &cli.App{
		Name:    "hanayo",
		Usage:   "The Ripple front end server",
		Version: hanayo.Version,
		Commands: []*cli.Command{
			{
				Name:   "web",
				Action: Web,
			},
		},
		Action: Web,
	}
	app.Run(os.Args)
}

// Web starts the Hanayo HTTP server.
func Web(ctx *cli.Context) error {
	srv := &http.Server{}
	err := srv.SetUpRoutes()
	if err != nil {
		return err
	}
	return nhttp.ListenAndServe(conf.Conf.ListenTo, srv)
}
