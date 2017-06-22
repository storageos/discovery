package store

import (
	"errors"
)

// KVPair represents the results of an operation on KVDB.
type KVPair struct {
	// Key for this kv pair.
	Key string
	// Value for this kv pair
	Value []byte

	// TTL value after which this key will expire from KVDB
	TTL int64
}

// KVFlags options for operations on KVDB
type KVFlags uint64

const (
	// KVPrevExists flag to check key already exists
	KVPrevExists KVFlags = 1 << iota
	// KVCreatedIndex flag compares with passed in index (possibly in KVPair)
	KVCreatedIndex
	// KVModifiedIndex flag compares with passed in index (possibly in KVPair)
	KVModifiedIndex
	// KVTTL uses TTL val from KVPair.
	KVTTL
)

// KVPairs list of KVPairs
type KVPairs []*KVPair

// Store - generic store interface
type Store interface {
	// Create is the same as Put except that ErrExist is returned if the key exists.
	Create(key string, value []byte, ttl int64) (*KVPair, error)

	Put(key string, value []byte, ttl int64) (*KVPair, error)

	// Get returns KVPair that maps to specified key or ErrNotFound.
	Get(key string) (*KVPair, error)
	// Delete deletes the KVPair specified by the key. ErrNotFound is returned
	// if the key is not found. The old KVPair is returned if successful.
	Delete(key string) error
}

var (
	// ErrNotFound raised if Key is not found
	ErrNotFound = errors.New("Key not found")

	// ErrExist raised if key already exists
	ErrExist = errors.New("Key already exists")

	// ErrValueMismatch raised if existing KVDB value mismatches with user provided value
	ErrValueMismatch = errors.New("Value mismatch")
)
