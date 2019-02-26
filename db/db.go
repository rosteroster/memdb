package db

import (
	"sync"
)

// Note: currently supports only one transaction at a time but allows multiple reads/lists

type MemDB struct {
	storage         map[string]string
	lock            sync.RWMutex
	transactionLock sync.Mutex
}

func New() *MemDB {
	return &MemDB{
		storage: map[string]string{},
		lock:    sync.RWMutex{},
	}
}

func (db *MemDB) write(key, val string) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.storage[key] = val
}

// Read is a safe read op that can be done outside of transactions
func (db *MemDB) Read(key string) string {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return db.storage[key]
}

// List is a safe op that can be done outside of transactions
func (db *MemDB) List() []string {
	db.lock.RLock()
	defer db.lock.RUnlock()

	list := []string{}
	for key := range db.storage {
		list = append(list, key)
	}
	return list
}

func (db *MemDB) delete(key string) {
	db.lock.Lock()
	defer db.lock.Unlock()

	delete(db.storage, key)
}

// Do is the only exported way to make changes to the DB (writes, deletes and rollbacks)
func (db *MemDB) Do(transaction *MemDBTransaction) {
	db.transactionLock.Lock()
	defer db.transactionLock.Unlock()

	for _, action := range transaction.Actions {
		action.do()
	}
}
