// -------------------------------------------------
// Package persistence
// Author: hanzhi
// Date: 2025/1/11
// -------------------------------------------------

package persistence

import (
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

// CGO_ENABLED=0 go get github.com/dgraph-io/badger

type Persistence interface {
	// 保存文档状态
	SaveState(docName string, state []byte) error
	// 加载文档状态
	LoadState(docName string) ([]byte, error)
	// 存储更新
	StoreUpdate(docName string, update []byte) error
	// 关闭存储
	Close() error
}

type BadgerPersistence struct {
	db *badger.DB
}

func NewBadgerPersistence(path string) (*BadgerPersistence, error) {
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerPersistence{db: db}, nil
}

func (p *BadgerPersistence) SaveState(docName string, state []byte) error {
	return p.db.Update(func(txn *badger.Txn) error {
		key := []byte("state:" + docName)
		entry := badger.NewEntry(key, state).WithTTL(24 * time.Hour)
		return txn.SetEntry(entry)
	})
}

func (p *BadgerPersistence) LoadState(docName string) ([]byte, error) {
	var state []byte
	err := p.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("state:" + docName))
		if err != nil {
			return err
		}
		state, err = item.ValueCopy(nil)
		return err
	})
	return state, err
}

func (p *BadgerPersistence) StoreUpdate(docName string, update []byte) error {
	return p.db.Update(func(txn *badger.Txn) error {
		key := []byte("update:" + docName + ":" + time.Now().Format(time.RFC3339Nano))
		return txn.Set(key, update)
	})
}

func (p *BadgerPersistence) Close() error {
	return p.db.Close()
}
