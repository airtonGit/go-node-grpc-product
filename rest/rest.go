package rest

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/peterbourgon/ff/v3"
)

//AppParams banco e listenAddr
type AppParams struct {
	ListenAddr string
}

//Setup flags de CLI e conexao banco
func Setup() (AppParams, error) {
	fs := flag.NewFlagSet("hospedare", flag.ExitOnError)
	var (
		listenAdr = fs.String("listen-addr", "127.0.0.1:8000", "Hospedare listen addr")
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("REST"),
		ff.WithConfigFile("product.conf"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		log.Println("FlagsFirst err:", err)
	}

	return AppParams{*listenAdr}, nil
}

//Listen inicia serviço
func Listen() error {
	params, err := Setup()
	if err != nil {
		return err
	}
	router := NewRouter(params)

	srv := &http.Server{
		Handler:      router,
		Addr:         params.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second}

	return srv.ListenAndServe()
}
