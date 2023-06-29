package main

import (
	"bookstore/internal/catalog/repo"
	"bookstore/internal/catalog/service"
	"bookstore/pkg/util"
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

	r.Route("/catalog", func(r chi.Router) {
		r.With(srv.Page).Get("/page/{i}", util.Dummy)
	})

	return r
}

func main() {
	booksRepo, err := repo.Connect("catalog", "books")
	if err != nil {
		log.Fatal(err)
	}
	defer booksRepo.Close()

	srv := (&service.Service{}).Init(booksRepo)

	go func() {
		err := http.ListenAndServe(
			":20",
			Run(srv),
		)
		log.Fatal(err)
	}()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
