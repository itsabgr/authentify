package authentify

import (
	"context"
	"github.com/google/uuid"
	"io"
	"time"
)
import badger "github.com/dgraph-io/badger/v3"

type Repo interface {
	Store(ctx context.Context, id string, val []byte, deadline time.Time) (err error)
	FindByID(ctx context.Context, id string) ([]byte, error)
	io.Closer
}
type badgerRepo struct {
	db *badger.DB
}

func (repo *badgerRepo) Store(ctx context.Context, id string, val []byte, deadline time.Time) (err error) {
	key := []byte(id)
	//
	tx := repo.db.NewTransaction(true)
	defer tx.Discard()
	item, err := tx.Get(key)
	if err == nil && false == item.IsDeletedOrExpired() {
		return badger.ErrConflict
	}
	err = tx.SetEntry(&badger.Entry{
		Key:       key,
		Value:     val,
		ExpiresAt: uint64(deadline.Unix()),
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	//
	return nil
}

func (repo *badgerRepo) FindByID(ctx context.Context, id string) ([]byte, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	key, err := uuid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	//
	tx := repo.db.NewTransaction(false)
	defer tx.Discard()
	item, err := tx.Get(key)
	if err != nil {
		return nil, err
	}
	if item.IsDeletedOrExpired() {
		return nil, badger.ErrKeyNotFound
	}
	return item.ValueCopy(nil)
}

func (repo *badgerRepo) Close() error {
	return repo.db.Sync()
}

func BadgerAsRepo(db *badger.DB) (Repo, error) {
	return &badgerRepo{
		db: db,
	}, nil
}
