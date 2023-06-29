package service

import (
	"bookstore/internal/auth/config"
	"bookstore/internal/auth/repo"
	"context"

	"github.com/go-chi/jwtauth/v5"
)

type Service struct {
	credentialsRepo repo.Repository
	sessionsRepo    repo.Repository
	tokenAuth       *jwtauth.JWTAuth
}

func (s *Service) Init(credentialsRepo, sessionsRepo repo.Repository) (self *Service) {
	s.credentialsRepo = credentialsRepo
	s.sessionsRepo = sessionsRepo
	s.tokenAuth = jwtauth.New("HS256", []byte(config.AuthSecret), nil)
	return s
}

func (s *Service) WatchExpiredSessions(ctx context.Context) error {
	return s.sessionsRepo.WatchExpirations(ctx)
}
