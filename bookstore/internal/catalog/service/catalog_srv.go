package service

import "bookstore/internal/catalog/repo"

type Service struct {
	booksRepo repo.Repository
}

func (s *Service) Init(booksRepo repo.Repository) (self *Service) {
	s.booksRepo = booksRepo
	return s
}
