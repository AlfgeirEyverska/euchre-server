package server

import (
	"net"
	"sync"
)

// ConnTracker keeps track of active connections and ensures proper closure and synchronization.
type ConnTracker struct {
	mu    sync.Mutex
	conns map[net.Conn]struct{}
	wg    sync.WaitGroup
}

func NewConnTracker() ConnTracker {
	return ConnTracker{
		mu:    sync.Mutex{},
		conns: make(map[net.Conn]struct{}),
		wg:    sync.WaitGroup{},
	}
}

// add registers a new connection and increments the WaitGroup counter.
func (ct *ConnTracker) add(conn net.Conn) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.conns[conn] = struct{}{}
	ct.wg.Add(1)
}

// done removes a connection and decrements the WaitGroup counter.
func (ct *ConnTracker) done(conn net.Conn) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if _, ok := ct.conns[conn]; ok {
		delete(ct.conns, conn)
		ct.wg.Done()
	}
}

// closeAll closes all tracked connections and removes them from the map.
func (ct *ConnTracker) CloseAll() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	for conn := range ct.conns {
		_ = conn.Close() // ignore error — shutting down anyway
		ct.wg.Done()
		delete(ct.conns, conn)
	}
}

func (ct *ConnTracker) Prune() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	for conn := range ct.conns {
		if !isAlive(conn) {
			_ = conn.Close() // ignore error — shutting down anyway
			ct.wg.Done()
			delete(ct.conns, conn)
		}
	}
}

// wait blocks until all tracked connections have completed.
func (ct *ConnTracker) Wait() {
	ct.wg.Wait()
}
