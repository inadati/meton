package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type TunnelRecipe struct {
	Auth        []ssh.AuthMethod
	GatewayUser string
	GatewayHost string
	DialAddr    string
	BindAddr    string
	Log         logger
}

type logger interface {
	Printf(string, ...interface{})
}

var Tunnel = &TunnelRecipe{}

func (t *TunnelRecipe) Forward(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sshClient, err := newReconnectableSSHClient(t.GatewayHost, &ssh.ClientConfig{
		User:            t.GatewayUser,
		Auth:            t.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("dial gateway %s: %w", t.GatewayHost, err)
	}
	defer sshClient.Close()
	go sshClient.KeepAlive(ctx)

	bindListener, err := closableListen(t.BindAddr)
	if err != nil {
		return fmt.Errorf("listen to bind address - %s: %w", t.BindAddr, err)
	}
	defer bindListener.Close()

	t.Log.Printf("start forwarding: %s -> %s", t.DialAddr, t.BindAddr)
	defer t.Log.Printf("stop forwarding: %s -> %s", t.DialAddr, t.BindAddr)

	t.startAccept(ctx, sshClient, bindListener)
	return nil
}

func (t *TunnelRecipe) startAccept(ctx context.Context, sshClient *ClientRecipe, bindListener *ListenerRecipe) {
	go func() {
		<-ctx.Done()
		bindListener.Close()
	}()

	for {
		bindConn, err := bindListener.Accept()
		if bindListener.IsClosed() {
			break
		}
		if err != nil {
			t.Log.Printf("failed to accept %s: %v", t.BindAddr, err)
			break
		}

		t.Log.Printf("accepted %s -> %s", t.BindAddr, bindConn.RemoteAddr())
		go func(bindConn net.Conn) {
			defer t.Log.Printf("disconnected %s -> %s", t.BindAddr, bindConn.RemoteAddr())

			connCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				<-connCtx.Done()
				bindConn.Close()
			}()

			dialConn, err := sshClient.Dial(ctx, "tcp", t.DialAddr)
			if err != nil {
				t.Log.Printf("failed to dial %s: %v", t.DialAddr, err)
				return
			}
			// ensure to close dial connection when copying finished.
			go func() {
				<-connCtx.Done()
				dialConn.Close()
			}()

			t.biCopy(dialConn, bindConn, cancel)
		}(bindConn)
	}
}

func (t *TunnelRecipe) biCopy(dialConn, bindConn net.Conn, shutdown func()) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := io.Copy(dialConn, bindConn); err != nil {
			t.Log.Printf("failed to copy %s -> %s: %v", t.DialAddr, t.BindAddr, err)
			shutdown()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := io.Copy(bindConn, dialConn); err != nil {
			t.Log.Printf("failed to copy %s -> %s: %v", t.BindAddr, t.DialAddr, err)
			shutdown()
		}
	}()

	wg.Wait()
}
