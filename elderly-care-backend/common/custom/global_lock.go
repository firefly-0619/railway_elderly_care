package custom

import (
	"sync"
	"time"
)

type GlobalLock struct {
	muLocks       sync.Mutex
	muRwLocks     sync.Mutex
	locks         map[string]*Lock
	rwLocks       map[string]*RWLock
	once          sync.Once
	cleanInterval time.Duration
}

type Lock struct {
	mutex sync.Mutex
	ttl   time.Duration // 锁的过期时间
}

type RWLock struct {
	rwMutex sync.RWMutex
	ttl     time.Duration // 锁的过期时间
}

func (g *GlobalLock) AcquireLock(key string, ttl time.Duration) *Lock {
	g.muLocks.Lock()
	lock, ok := g.locks[key]
	g.muLocks.Unlock()
	if !ok {
		g.muLocks.Lock()
		lock = &Lock{
			ttl: ttl,
		}
		g.locks[key] = lock
		g.muLocks.Unlock()
	} else {

	}
	return lock
}

func (g *GlobalLock) AcquireRwLock(key string, ttl time.Duration) *RWLock {
	g.muRwLocks.Lock()
	rwLock, ok := g.rwLocks[key]
	g.muRwLocks.Unlock()
	if !ok {
		g.muRwLocks.Lock()
		rwLock = &RWLock{
			ttl: ttl,
		}
		g.rwLocks[key] = rwLock
		g.muRwLocks.Unlock()
	}
	return rwLock
}

func (l *Lock) Lock() {
	if time.Now().UnixMilli() > l.ttl.Milliseconds() {
		l.mutex.Lock()
	}
}
