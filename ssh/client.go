package ssh

import (
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
)

type ClientRecipe struct {
	Host   string
	Config *ssh.ClientConfig

	Mux sync.RWMutex
	C   *ssh.Client
}

func newReconnectableSSHClient(host string, config *ssh.ClientConfig) (*ClientRecipe, error) {
	c, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}
	return &ClientRecipe{Host: host, Config: config, C: c}, nil
}

func (r *ClientRecipe) KeepAlive(ctx context.Context) {
	wait := make(chan error, 1)
	go func() {
		wait <- r.getC().Wait()
	}()

	var aliveErrCount uint32
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-wait:
			return
		case <-ticker.C:
			if atomic.LoadUint32(&aliveErrCount) > 1 {
				log.Printf("failed to keep alive of %v", r.getC().RemoteAddr())
				r.getC().Close()
				return
			}
		case <-ctx.Done():
			return
		}

		go func() {
			_, _, err := r.getC().SendRequest("keepalive@openssh.com", true, nil)
			if err != nil {
				atomic.AddUint32(&aliveErrCount, 1)
			}
		}()
	}
}

func (r *ClientRecipe) Close() error {
	return r.getC().Close()
}

func (r *ClientRecipe) getC() *ssh.Client {
	r.Mux.RLock()
	defer r.Mux.RUnlock()
	return r.C
}

func (r *ClientRecipe) setC(client *ssh.Client) {
	r.Mux.Lock()
	defer r.Mux.Unlock()
	r.C = client
}

func (r *ClientRecipe) Dial(ctx context.Context, n, addr string) (net.Conn, error) {
	conn, err := r.getC().Dial(n, addr)
	if err != nil {
		if rErr := r.reconnect(ctx); rErr != nil {
			return nil, err
		}
		return r.getC().Dial(n, addr)
	}
	return conn, nil
}

func (r *ClientRecipe) reconnect(ctx context.Context) error {
	client, err := ssh.Dial("tcp", r.Host, r.Config)
	if err != nil {
		return err
	}
	r.getC().Close()
	r.setC(client)
	go r.KeepAlive(ctx)
	return nil
}