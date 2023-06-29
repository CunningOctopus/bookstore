package main

import (
	"bookstore/internal/auth/repo"
	"bookstore/internal/auth/service"
	"bookstore/pkg/util"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Run(srv *service.Service) http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.With(srv.Signup).Post("/signup", util.Dummy)
		r.With(srv.Signin).Post("/signin", util.Dummy)
		r.With(srv.Logout).Post("/logout", util.Dummy)
	})

	return r
}

func main() {
	credentialsRepo, err := repo.Connect("auth", "credentials")
	if err != nil {
		log.Fatal(err)
	}
	defer log.Println(credentialsRepo.Close())

	sessionsRepo, err := repo.Connect("auth", "sessions")
	if err != nil {
		log.Fatal(err)
	}
	defer log.Println(sessionsRepo.Close())

	srv := (&service.Service{}).Init(credentialsRepo, sessionsRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go srv.WatchExpiredSessions(ctx)

	go func() {
		err := http.ListenAndServe(
			":10",
			Run(srv),
		)
		if err != nil {
			log.Print(err)
			cancel()
			os.Exit(1)
		}
	}()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
