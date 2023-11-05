package datastore

import (
	"fmt"
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

func newElem() elem {
	return elem{
		stamp:       time.Now(),
		subscribers: make([]Callback, 0),
		v:           nil,
	}
}

func New() *DataStore {
	return &DataStore{
		store: make(map[string]elem),
		lock:  sync.RWMutex{},
	}
}

func (d *DataStore) Publish(key string, v interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()

	e, ok := d.store[key]
	if !ok {
		e = newElem()
	}

	e.v = v
	e.stamp = time.Now()

	d.store[key] = e

	for _, f := range e.subscribers {
		f(key, e.stamp, e.v)
	}
}

func (d *DataStore) Get(key string) (interface{}, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	e, ok := d.store[key]
	if !ok {
		return nil, fmt.Errorf("%s: not found", key)
	}

	return e.v, nil
}

func (d *DataStore) Register(keys []string, f Callback) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, k := range keys {
		e, ok := d.store[k]
		if !ok {
			e = newElem()
		}
		e.subscribers = append(e.subscribers, f)
		d.store[k] = e
	}
}
