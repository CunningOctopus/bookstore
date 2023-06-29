package service

import (
	"bookstore/internal/auth/config"
	"bookstore/internal/auth/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
curl \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=some@email.world&password=passw0rd" \
  http://localhost:8080/signin
*/

func (srv *Service) Signin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		log.Println(email, password)

		hashedPassword := srv.credentialsRepo.Get([]byte(email))
		err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			log.Printf("attempt to login with incorrect credentials: %s", err)
			return
		}

		sessID := uuid.NewString()
		expires := time.Now().Add(config.TokenExpTime).UTC()
		session := domain.Session{
			EMail:   email,
			Expires: expires,
		}
		bytes, err := json.Marshal(session)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to marhsal session data: %s", err)
			return
		}

		if err := srv.sessionsRepo.Insert([]byte(sessID), bytes); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to store session data: %s", err)
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:    "session",
			Value:   sessID,
			Expires: expires,
		})

		log.Println(string(email), string(hashedPassword))

		fmt.Fprint(rw, "signin successful")

		next.ServeHTTP(rw, r)
	})
}
