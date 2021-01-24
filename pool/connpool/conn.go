package connpool

import (
	"net"
	"sync"
	"time"
)

type Conn struct {
	net.Conn
	c *channelPool
	unusable bool		// if unusable is true, the conn should be closed
	mu sync.RWMutex
	t time.Time  // connection idle time
	dialTimeout time.Duration // connection timeout duration
}

func (p *Conn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

// 假关闭,重新放回连接池进行复用
func (p *Conn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
	}

	// reset connection deadline
	p.Conn.SetDeadline(time.Time{})

	return p.c.Put(p)
}

func (p *Conn) Read(b []byte) (int, error) {
	if p.unusable {
		return 0, ErrConnClosed
	}
	n, err := p.Conn.Read(b)
	if err != nil {
		p.MarkUnusable()
		p.Conn.Close()
	}
	return n, err
}

func (p *Conn) Write(b []byte) (int, error) {
	if p.unusable {
		return 0, ErrConnClosed
	}
	n, err := p.Conn.Write(b)
	if err != nil {
		p.MarkUnusable()
		p.Conn.Close()
	}
	return n, err
}
