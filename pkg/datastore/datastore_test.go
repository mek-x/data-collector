package datastore

import (
	"reflect"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *DataStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataStore_Publish(t *testing.T) {
	type fields struct {
		store map[string]elem
		lock  sync.RWMutex
	}
	type args struct {
		key string
		v   interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataStore{
				store: tt.fields.store,
				lock:  tt.fields.lock,
			}
			d.Publish(tt.args.key, tt.args.v)
		})
	}
}

func TestDataStore_Get(t *testing.T) {
	type fields struct {
		store map[string]elem
		lock  sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataStore{
				store: tt.fields.store,
				lock:  tt.fields.lock,
			}
			got, err := d.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataStore.Get(%v) error = %v, wantErr %v", tt.args.key, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataStore.Get(%v) = %v, want %v", tt.args.key, got, tt.want)
			}
		})
	}
}

func TestDataStore_Register(t *testing.T) {
	type fields struct {
		store map[string]elem
		lock  sync.RWMutex
	}
	type args struct {
		keys []string
		f    Callback
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataStore{
				store: tt.fields.store,
				lock:  tt.fields.lock,
			}
			d.Register(tt.args.keys, tt.args.f)
		})
	}
}
