package repo

import (
	"bookstore/internal/catalog/domain"
	"encoding/json"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

var ErrValueDoesntExist = errors.New("value at key doesn't exist")
var ErrRangeOutsideKeyN = errors.New("range is outside bucket's KeyN")

type RepoConn struct {
	conn   *bbolt.DB
	db     string
	bucket string
}

func Connect(db, bucket string) (*RepoConn, error) {

	conn, err := bbolt.Open(db+".db", 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %w", db, err)
	}

	conn.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
		}
		return nil
	})

	return &RepoConn{
		conn:   conn,
		db:     db,
		bucket: bucket,
	}, nil
}

func (repo *RepoConn) Close() error {

	if err := repo.conn.Close(); err != nil {
		return err
	}
	repo.db = ""
	repo.bucket = ""
	return nil
}

func (repo *RepoConn) DB() string { return repo.db }

func (repo *RepoConn) Bucket() string { return repo.bucket }

// Inserts `v` at `k` if `k` is unique. Otherwise, returns error.
func (repo *RepoConn) Insert(k []byte, v domain.Book) error {
	return repo.conn.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(repo.bucket))

		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}

		return bkt.Put(k, bytes)
	})
}

func (repo *RepoConn) Get(k []byte) (v []byte) {
	repo.conn.View(func(tx *bbolt.Tx) error {
		v = tx.Bucket([]byte(repo.bucket)).Get(k)
		return nil
	})

	return v
}

// Gets range of records. Records are sorted by key in lexicographical order.
// If number of pairs in the bucket less then `left`, returns empty slice.
// If number of pairs in the bucket less then `right`, reduces `right` down to number of pairs.
// Left included, right excluded.
func (repo *RepoConn) GetRange(left, right int) ([][]byte, error) {
	if left < 0 || left > right {
		return [][]byte{}, ErrRangeOutsideKeyN
	}

	var r [][]byte
	err := repo.conn.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(repo.bucket))

		keyn := bkt.Stats().KeyN
		if keyn < left {
			return ErrRangeOutsideKeyN
		}
		if keyn < right {
			right = keyn
		}

		c := bkt.Cursor()
		for cart := 0; cart != left; cart++ {
			c.Next()
		}

		r = make([][]byte, right-left)
		for i := range r {
			_, r[i] = c.Next()
		}

		return nil
	})

	return r, err
}

// Updates value at `k` with a `new`.
func (repo *RepoConn) Update(k []byte, new domain.Book) error {
	return repo.conn.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(repo.bucket))

		bytes, err := json.Marshal(new)
		if err != nil {
			return err
		}

		return bkt.Put(k, bytes)
	})
}

func (repo *RepoConn) Delete(k []byte) error {
	return repo.conn.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(repo.bucket))

		if v := bkt.Get(k); v == nil {
			return ErrValueDoesntExist
		}

		if err := bkt.Delete(k); err != nil {
			return fmt.Errorf("failed to delete: %s", err)
		}

		return nil
	})
}
