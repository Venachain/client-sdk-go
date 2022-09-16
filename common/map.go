package common

import (
	"sync"

	common_venachain "github.com/Venachain/client-sdk-go/venachain/common"
)

type HashMaps struct {
	hashMaps []hashMap
}

func NewHashMaps() HashMaps {
	maps := make([]hashMap, 16)
	for i := 0; i < 16; i++ {
		maps[i] = hashMap{m: make(map[common_venachain.Hash]chan struct{})}
	}
	return HashMaps{
		hashMaps: maps,
	}
}

func (h HashMaps) Put(hash common_venachain.Hash, ch chan struct{}) {
	h.hashMaps[getIndex(hash)].put(hash, ch)
}

func (h HashMaps) Delete(hash common_venachain.Hash) {
	h.hashMaps[getIndex(hash)].delete(hash)
}

func (h HashMaps) Contains(hash common_venachain.Hash) bool {
	return h.hashMaps[getIndex(hash)].contains(hash)
}

func getIndex(hash common_venachain.Hash) int {
	v := hash[len(hash)-1]
	return int(v % 16)
}

type hashMap struct {
	m    map[common_venachain.Hash]chan struct{}
	lock sync.RWMutex
}

func (t hashMap) contains(hash common_venachain.Hash) bool {
	t.lock.RLock()
	ch, ok := t.m[hash]
	t.lock.RUnlock()
	if ok {
		ch <- struct{}{}
		t.lock.Lock()
		defer t.lock.Unlock()
		delete(t.m, hash)
	}
	return ok
}

func (t hashMap) put(hash common_venachain.Hash, ch chan struct{}) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.m[hash] = ch
}

func (t hashMap) delete(hash common_venachain.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.m, hash)
}
