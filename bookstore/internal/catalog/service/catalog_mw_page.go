package service

import (
	"bookstore/internal/catalog/config"
	"bookstore/internal/catalog/repo"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (srv *Service) Page(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		pageStr := chi.URLParam(r, "i")
		if pageStr == "" {
			rw.WriteHeader(http.StatusBadRequest)
			log.Println("no page number was supplied in the URL")
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		left := config.ItemsOnPage * page
		right := config.ItemsOnPage*page + config.ItemsOnPage
		books, err := srv.booksRepo.GetRange(left, right)
		if err != nil {
			if err == repo.ErrRangeOutsideKeyN {
				rw.WriteHeader(http.StatusBadRequest)
			}
			rw.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte("["))
		for _, b := range books {
			rw.Write(b)
		}
		rw.Write([]byte("]"))

		next.ServeHTTP(rw, r)
	})
}
