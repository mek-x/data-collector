package datastore

import (
	"fmt"
	"sync"
	"time"
)

type Callback func(key string, t time.Time, v, old any)
type TimeoutF func(key string, t time.Time, v any)

type elem struct {
	stamp       time.Time
	v           interface{}
	subscribers []Callback
	timer       *time.Timer
	timeout     time.Duration
	tf          TimeoutF
}

type DataStore interface {
	Publish(key string, v interface{})
	Get(key string) (interface{}, time.Time, error)
	Register(keys []string, f Callback, timeout time.Duration, tf TimeoutF)
	Map() map[string]any
}

type dataStore struct {
	store map[string]elem
	lock  sync.RWMutex
}

var _ DataStore = (*dataStore)(nil)

func newElem() elem {
	return elem{
		stamp:       time.Now(),
		subscribers: make([]Callback, 0),
		v:           nil,
		timer:       nil,
		timeout:     0,
		tf:          nil,
	}
}

func New() *dataStore {
	return &dataStore{
		store: make(map[string]elem),
		lock:  sync.RWMutex{},
	}
}

func (d *dataStore) Publish(key string, v interface{}) {
	d.lock.Lock()

	e, ok := d.store[key]
	if !ok {
		e = newElem()
	} else {
		if e.timer != nil && e.timeout > 0 {
			e.timer.Reset(e.timeout)
		}
	}

	old := e.v
	e.v = v
	e.stamp = time.Now()

	d.store[key] = e

	d.lock.Unlock()

	for _, f := range e.subscribers {
		f(key, e.stamp, v, old)
	}
}

func (d *dataStore) Get(key string) (interface{}, time.Time, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	e, ok := d.store[key]
	if !ok {
		return nil, time.Time{}, fmt.Errorf("%s: not found", key)
	}

	return e.v, e.stamp, nil
}

func (d *dataStore) Register(keys []string, f Callback, timeout time.Duration, tf TimeoutF) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, k := range keys {
		e, ok := d.store[k]
		if !ok {
			e = newElem()
		} else {
			if e.timer != nil && e.timeout > 0 {
				e.timer.Stop()
			}
		}

		e.subscribers = append(e.subscribers, f)

		if timeout > 0 {
			e.timeout = timeout
			e.tf = tf
			e.timer = time.AfterFunc(e.timeout, func() {
				e.tf(k, e.stamp, e.v)
			})
		}
		d.store[k] = e
	}
}

func (d *dataStore) Map() map[string]any {
	out := make(map[string]any)

	d.lock.Lock()
	defer d.lock.Unlock()

	for k, v := range d.store {
		out[k] = v.v
	}

	return out
}
