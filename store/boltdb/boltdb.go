package boltdb

import (
	"bytes"
	"sync"
	"time"

	"github.com/storageos/discovery/store"
	// "github.com/storageos/discovery/util/codecs"

	"github.com/boltdb/bolt"
)

type Store struct {
	db               *bolt.DB
	mu               *sync.Mutex
	index            uint64
	tokensBucketName []byte
}

func New(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	st := &Store{
		db:               db,
		mu:               &sync.Mutex{},
		tokensBucketName: []byte("tokens"),
	}
	// ensure bucket
	err = st.ensureBuckets()
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Store) ensureBuckets() error {
	// store some data
	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.tokensBucketName)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Create(key string, value []byte, ttl int64) (*store.KVPair, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.Get(key)
	if err != nil {
		return s.put(key, value, ttl)
	}

	return nil, store.ErrExist
}

func (s *Store) Put(key string, value []byte, ttl int64) (*store.KVPair, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.put(key, value, ttl)
}

func (s *Store) Get(key string) (*store.KVPair, error) {
	buf := bytes.Buffer{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.tokensBucketName)
		v := b.Get([]byte(key))
		if v == nil {
			return store.ErrNotFound
		}
		buf.Write(v)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &store.KVPair{
		Key:   key,
		Value: buf.Bytes(),
	}, nil
}

func (s *Store) put(key string, value []byte, ttl int64) (*store.KVPair, error) {
	var err error

	if ttl != 0 {
		time.AfterFunc(time.Second*time.Duration(ttl), func() {
			s.Delete(key)
		})
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.tokensBucketName)
		err := b.Put([]byte(key), value)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &store.KVPair{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}, nil
}

func (s *Store) Delete(key string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.tokensBucketName).Delete([]byte(key))
	})

	return err
}
