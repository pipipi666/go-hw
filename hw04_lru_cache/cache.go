package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type cacheItemValue struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	val := cacheItemValue{key, value}
	item, ok := l.items[key]

	if ok {
		item.Value = val
		l.queue.MoveToFront(item)

		return true
	}

	newItem := l.queue.PushFront(val)
	l.items[key] = newItem

	if l.queue.Len() > l.capacity {
		tail := l.queue.Back()
		v, ok := tail.Value.(cacheItemValue)

		if ok {
			delete(l.items, v.key)
		}

		l.queue.Remove(tail)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item)
		v, ok := item.Value.(cacheItemValue)

		if ok {
			return v.value, true
		}

		return nil, false
	}

	return nil, false
}

func (l *lruCache) Clear() {
	for k, v := range l.items {
		l.mu.Lock()

		delete(l.items, k)
		l.queue.Remove(v)

		l.mu.Unlock()
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
