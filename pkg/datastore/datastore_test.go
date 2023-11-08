package datastore

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestDataStore_Publish(t *testing.T) {
	type fields struct {
		store map[string]elem
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
			d := &dataStore{
				store: tt.fields.store,
				lock:  sync.RWMutex{},
			}
			d.Publish(tt.args.key, tt.args.v)
		})
	}
}

func TestDataStore_Get(t *testing.T) {
	type fields struct {
		store map[string]elem
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
		{
			name:    "simple",
			fields:  fields{store: map[string]elem{"a": {v: 1}}},
			args:    args{"a"},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dataStore{
				store: tt.fields.store,
				lock:  sync.RWMutex{},
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

func dummyCallback(k string, t time.Time, v interface{}) {}

func TestDataStore_Register(t *testing.T) {
	type args struct {
		keys []string
		f    Callback
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple",
			args: args{
				keys: []string{"a"},
				f:    dummyCallback,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			d.Register(tt.args.keys, tt.args.f)
			for _, v := range tt.args.keys {
				for _, s := range d.store[v].subscribers {
					if s != nil {
						goto found
					}
				}
			}
			t.Errorf("DataStore.Register(%v,%v) subsciber not found in the store", tt.args.keys, tt.args.f)
		found:
		})
	}
}
