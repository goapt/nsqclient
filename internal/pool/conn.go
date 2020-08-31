package pool

import (
	"sync"

	nsq "github.com/nsqio/go-nsq"
)

// Producer is a wrapper around net.Conn to modify the the behavior of
// net.Conn's Close() method.
type Producer struct {
	*nsq.Producer
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

// Close puts the given connects back to the pool instead of closing it.
func (p *Producer) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Producer != nil {
			p.Producer.Stop()
			return nil
		}
		return nil
	}
	return p.c.put(p.Producer)
}

// MarkUnusable marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (p *Producer) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

// newConn wraps a standard net.Conn to a poolConn net.Conn.
func (c *channelPool) wrapConn(conn *nsq.Producer) *Producer {
	p := &Producer{c: c}
	p.Producer = conn
	return p
}
