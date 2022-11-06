package keysutil

import lru "github.com/hashicorp/golang-lru"

type TransitLRU struct {
	size int
	lru  *lru.TwoQueueCache
}

func NewTransitLRU(size int) (*TransitLRU, error) {
	lru, err := lru.New2Q(size)
	return &TransitLRU{lru: lru, size: size}, err
}

func (c *TransitLRU) Delete(key any) {
	c.lru.Remove(key)
}

func (c *TransitLRU) Load(key any) (value any, ok bool) {
	return c.lru.Get(key)
}

func (c *TransitLRU) Store(key, value any) {
	c.lru.Add(key, value)
}

func (c *TransitLRU) Size() int {
	return c.size
}
