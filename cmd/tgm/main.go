package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/xerrors"
)

func run(ctx context.Context) error {
	app, err := NewApp()
	if err != nil {
		return xerrors.Errorf("config: %w", err)
	}

	defer func() {
		err = app.Close()

		if err != nil {
			log.Println(app.Close())
		}
	}()

	app.router.HandleFunc("/messages", func(rw http.ResponseWriter, r *http.Request) {
		messages, err := app.msgRepo.GetAll()

		if err != nil {
			err = json.NewEncoder(rw).Encode(err)
		} else {
			err = json.NewEncoder(rw).Encode(&messages)
		}

		if err != nil {
			log.Println(err)
		}
	}).Methods("GET")

	return app.Run(ctx)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Println(err)
	}
}
