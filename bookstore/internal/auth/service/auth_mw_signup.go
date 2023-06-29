package service

import (
	"bookstore/internal/auth/config"
	"bookstore/internal/auth/domain"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

/*
curl \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=some@email.world&password=passw0rd" \
  http://localhost:8080/signup
*/

func (srv *Service) Signup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		creds := domain.Credentials{
			EMail:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}
		if err := creds.Validate(); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		// TODO: Hash EMail as well
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), config.HashingCost)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := srv.credentialsRepo.Insert([]byte(creds.EMail), hashedPassword); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		log.Println(string(creds.EMail), string(hashedPassword))
		log.Println(creds)

		next.ServeHTTP(rw, r)
	})
}
