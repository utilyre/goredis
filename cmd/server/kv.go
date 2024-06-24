package main

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type KV struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKV() *KV {
	return &KV{data: map[string][]byte{}}
}

func (kv *KV) Set(key string, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[key] = val
	return nil
}

func (kv *KV) Get(key string) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	return val, nil
}
