package repo

import (
	"bookstore/internal/auth/config"
	"bookstore/internal/auth/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"go.etcd.io/bbolt"
)

var ErrValueDoesntExist = errors.New("value at key doesn't exist")
var ErrNotUnique = errors.New("key is not unique")
var ErrBucketDoesntExist = errors.New("bucket doesn't exist")

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
func (repo *RepoConn) Insert(k, v []byte) error {
	return repo.conn.Update(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket([]byte(repo.bucket))

		if v := bkt.Get(k); v != nil {
			return fmt.Errorf("%w: %s", ErrNotUnique, string(k))
		}

		err := bkt.Put(k, v)
		if err != nil {
			return fmt.Errorf("failed to put %s into bucket %s: %w",
				string(v),
				repo.bucket,
				err,
			)
		}

		return nil
	})
}

func (repo *RepoConn) Get(k []byte) (v []byte) {
	repo.conn.View(func(tx *bbolt.Tx) error {
		v = tx.Bucket([]byte(repo.bucket)).Get(k)
		return nil
	})

	return
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

func (repo *RepoConn) RemoveExpired() error {

	tx, err := repo.conn.Begin(true)
	if err != nil {
		return fmt.Errorf("failed to initiate transaction: %w", err)
	}

	bkt := tx.Bucket([]byte(repo.bucket))
	c := bkt.Cursor()
	now := time.Now()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		s := domain.Session{}
		if err := json.Unmarshal(v, &s); err != nil {
			return err
		}

		if s.Expires.Before(now) {
			/* ErrIncompatibleValue may occur if you have nested bucket
			which may become the case in the future. */
			err = c.Delete()
			if err != nil {
				return err
			}
			log.Printf("%s: session expired\n", string(k))
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (repo *RepoConn) WatchExpirations(ctx context.Context) error {
	timeout := time.After(time.Second * 0) // expired records are deleted immediately on startup

	for {
		select {
		case <-timeout:
			if err := repo.RemoveExpired(); err != nil {
				return err
			}
			timeout = time.After(time.Second * config.RepoChangestreamSleepSeconds)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
