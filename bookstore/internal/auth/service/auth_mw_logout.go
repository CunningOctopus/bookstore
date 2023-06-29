package service

import (
	"bookstore/internal/auth/repo"
	"log"
	"net/http"
)

func (srv *Service) Logout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("no session cookie was supplied")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := srv.sessionsRepo.Delete([]byte(cookie.Value)); err != nil {

			if err == repo.ErrValueDoesntExist {
				log.Println("invalid session cookie was supplied")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
