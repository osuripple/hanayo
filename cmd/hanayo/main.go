package main

import (
	"fmt"
	nhttp "net/http"
	"os"
	"time"

	"git.zxq.co/ripple/hanayo"
	"git.zxq.co/ripple/hanayo/helpers/conf"
	"git.zxq.co/ripple/hanayo/http"
	"git.zxq.co/ripple/hanayo/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/urfave/cli.v2"
)

var start = time.Now()

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
	db, err := sqlx.Open("mysql", conf.Conf.DSN)
	if err != nil {
		return err
	}
	sp := &mysql.ServiceProvider{db}

	srv := &http.Server{
		UserService: sp,
		TFAService:  sp,
	}
	err = srv.SetUpSimplePages()
	if err != nil {
		return err
	}
	srv.SetUpFuncMap()
	fmt.Println("Starting listening on", conf.Conf.ListenTo, "- startup took",
		time.Since(start))
	return nhttp.ListenAndServe(conf.Conf.ListenTo, srv)
}
