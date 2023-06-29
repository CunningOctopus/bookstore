package repo

import "context"

type Repository interface {
	Close() error

	DB() string     // db msut be private
	Bucket() string // buket must be private

	Insert(k, v []byte) error
	Get(k []byte) (v []byte)
	Delete(k []byte) error

	WatchExpirations(ctx context.Context) error
}
