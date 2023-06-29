package repo

import "bookstore/internal/catalog/domain"

type Repository interface {
	Close() error

	DB() string     // db must be private
	Bucket() string // buket must be private

	Insert(k []byte, v domain.Book) error
	Get(k []byte) (v []byte)
	GetRange(left, right int) ([][]byte, error)
	Update(k []byte, new domain.Book) error
	Delete(k []byte) error
}
