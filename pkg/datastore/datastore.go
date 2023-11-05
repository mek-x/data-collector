package datastore

import (
	"sync"
	"time"
)

type Callback func(key string, t time.Time, v interface{})

type elem struct {
	stamp       time.Time
	v           interface{}
	subscribers []Callback
}

type DataStore struct {
	store map[string]elem
	lock  sync.RWMutex
}

func New() *DataStore {
	return &DataStore{
		store: make(map[string]elem),
		lock:  sync.RWMutex{},
	}
}

func (d *DataStore) Publish(key string, v interface{}) error {
	return nil
}

func (d *DataStore) Get(key string) (interface{}, error) {
	return nil, nil
}

func (d *DataStore) Register(keys string, f Callback) error {
	return nil
}
