package rpc

import (
	"context"
	"sync/atomic"

	"gitlab.com/vocdoni/go-dvote/log"
	"nhooyr.io/websocket"
)

const (
	free = 0
	busy = 1
)

// TendermintRPC holds a pool of connections and keeps track of which are available
type TendermintRPC struct {
	Conns []PoolConnection
	// Index of latest-used connection
	roundRobinIndex uint32
}

// AddConnection adds a connection to the pool and initializes it as available
func (t *TendermintRPC) AddConnection(c *websocket.Conn) {
	t.Conns = append(t.Conns, PoolConnection{
		C:         c,
		available: free,
	})
}

// GetConnection finds, locks, and returns the next available poolconnection. Caller is responsible for releasing the connection.
func (t *TendermintRPC) GetConnection() *PoolConnection {
	for i := int(atomic.LoadUint32(&t.roundRobinIndex)); ; i++ {
		// If at the end of conns, round-robin back to beginning
		if i >= len(t.Conns) {
			i = 0
		}

		// Non-thread-safe check status: faster check if resource is NOT available, then move on
		if i < len(t.Conns) && t.Conns[i].Lock() {
			// Store returned index to start looking for next connection, so that we can search from the last claimed connection rather than the first connection every time. That would result in clustering around the beginning of the array and not using the later connections.
			atomic.StoreUint32(&t.roundRobinIndex, uint32(i))
			return &t.Conns[i]
		}
	}
}

// Close closes all connections in the pool
func (t *TendermintRPC) Close() {
	for _, conn := range t.Conns {
		conn.Close()
	}
	log.Infof("Closed %d websocket connections", len(t.Conns))
}

// PoolConnection holds a single websockets connection and a status int
type PoolConnection struct {
	C         *websocket.Conn
	available int32
}

// Close safely closes the connection
func (p *PoolConnection) Close() {
	atomic.StoreInt32(&p.available, busy)
	p.C.Close(websocket.StatusNormalClosure, "closed by caller")
}

// Lock returns false if the connection is already locked
func (p *PoolConnection) Lock() bool {
	return atomic.CompareAndSwapInt32(&p.available, free, busy)
}

// Release does not ensure that the connection is unavailable but sets it to be available either way.
func (p *PoolConnection) Release() {
	atomic.StoreInt32(&p.available, free)
}

// WriteRead executes a write & read operation on the websocket connection
func (p *PoolConnection) WriteRead(ctx context.Context, request []byte) ([]byte, error) {
	err := p.C.Write(ctx, websocket.MessageText, request)
	if err != nil {
		return nil, err
	}
	_, msg, err := p.C.Read(ctx)
	return msg, err
}
