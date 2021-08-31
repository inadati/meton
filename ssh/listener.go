package ssh

import (
	"net"
	"os"
	"regexp"
	"sync"
)

type ListenerRecipe struct {
	l net.Listener

	mux    sync.RWMutex
	closed bool
}

var (
	tcpAddressPattern = regexp.MustCompile(`(.+\.)+\w+:\d+`)
)

func closableListen(address string) (*ListenerRecipe, error) {
	var (
		l   net.Listener
		err error
	)
	if tcpAddressPattern.MatchString(address) {
		l, err = net.Listen("tcp", address)
	} else {
		// try unix socket connection
		// remove sock file is already exists
		if _, err := os.Stat(address); err == nil {
			_ = os.Remove(address)
		}
		l, err = net.Listen("unix", address)
	}
	if err != nil {
		return nil, err
	}
	return &ListenerRecipe{l: l}, nil
}

func (r *ListenerRecipe) Close() error {
	r.mux.Lock()
	r.closed = true
	r.mux.Unlock()
	return r.l.Close()
}

func (r *ListenerRecipe) IsClosed() bool {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.closed
}

func (r *ListenerRecipe) Accept() (net.Conn, error) {
	return r.l.Accept()
}