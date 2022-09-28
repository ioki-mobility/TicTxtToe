package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hashicorp/go-hclog"
	"github.com/ioki-mobility/TicTxtToe/backend"
	"github.com/ioki-mobility/TicTxtToe/frontend"
	"github.com/twilio/twilio-go"
)

type config struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phonenumber"`
}

func main() {
	lg := hclog.Default()
	ctx, can := context.WithCancel(context.Background())

	f, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	cfg := &config{}
	err = json.Unmarshal(b, cfg)
	if err != nil {
		panic(err)
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.Username,
		Password: cfg.Password,
	})

	twilioPhoneNumber := cfg.PhoneNumber

	bk := backend.NewBackend()
	fr := frontend.New(client, bk, twilioPhoneNumber, lg)
	go watchSignals(can, lg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/sms", func(w http.ResponseWriter, r *http.Request) {
		fr.ServeHTTP(w, r)
	})

	srv := http.Server{Addr: ":4040", Handler: r}

	go func() {
		_ = srv.ListenAndServe()
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
}

// watchSignals waits for one of the registered OS signals. Once a signal is received it calls
// fn and exits.
func watchSignals(can context.CancelFunc, lg hclog.Logger) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	lg.Info("received OS signal", "signal", <-sig)
	can()
}
