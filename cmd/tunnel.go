package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/meton888/meton/endpoint"
	"github.com/meton888/meton/ssh"
	"github.com/meton888/meton/syscallhelper"
	"github.com/urfave/cli/v2"
)

var TunnelCommand = &cli.Command{
	Name:  "tunnel",
	Usage: "Build port forwarding to Mesos Master, Marathon, Chronos.",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {

		var rLimit syscall.Rlimit
		if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
			log.Fatalf("failed to set ulimit: %v", err)
		}
		newRLimit := syscall.Rlimit{
			Cur: syscallhelper.RlimitMax(rLimit),
			Max: syscallhelper.RlimitMax(rLimit),
		}
		if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &newRLimit); err != nil {
			log.Fatalf("failed to set ulimit: %v", err)
		}

		gatewayAddr := endpoint.PrimaryNode.SSH(cfg.Cluster.Nodes.Master)
		mesosMasterRemoteHost := endpoint.PrimaryNode.MesosMaster(cfg.Cluster.Nodes.Master)
		marathonRemoteHost := endpoint.PrimaryNode.Marathon(cfg.Cluster.Nodes.Master)
		chronosRemoteHost := endpoint.PrimaryNode.Chronos(cfg.Cluster.Nodes.Master)

		logger := log.New(os.Stdout, "[sshtunnel] ", log.Flags())

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		defer stop()

		var wg sync.WaitGroup

		wg.Add(1)
		go func(kfp, gatewayStr string, remoteHost string, localHost string) {
			defer wg.Done()
			auth, _ := ssh.ParseKeyFile(kfp)
			ssh.Tunnel = &ssh.TunnelRecipe{
				Auth:        auth,
				GatewayUser: cfg.Cluster.Owner,
				GatewayHost: gatewayStr,
				DialAddr:    remoteHost,
				BindAddr:    localHost,
				Log:         logger,
			}
			if err := ssh.Tunnel.Forward(ctx); err != nil {
				log.Printf("failed to forward tunnel - %s -> %s: %v", remoteHost, localHost, err)
				stop()
			}
		}(cfg.Cluster.KeyFile, gatewayAddr, mesosMasterRemoteHost, endpoint.LocalHost.MesosMaster())
		
		wg.Add(1)
		go func(kfp, gatewayStr string, remoteHost string, localHost string) {
			defer wg.Done()
			auth, _ := ssh.ParseKeyFile(kfp)
			ssh.Tunnel = &ssh.TunnelRecipe{
				Auth:        auth,
				GatewayUser: cfg.Cluster.Owner,
				GatewayHost: gatewayStr,
				DialAddr:    remoteHost,
				BindAddr:    localHost,
				Log:         logger,
			}
			if err := ssh.Tunnel.Forward(ctx); err != nil {
				log.Printf("failed to forward tunnel - %s -> %s: %v", remoteHost, localHost, err)
				stop()
			}
		}(cfg.Cluster.KeyFile, gatewayAddr, marathonRemoteHost, endpoint.LocalHost.Marathon())
		
		wg.Add(1)
		go func(kfp, gatewayStr string, remoteHost string, localHost string) {
			defer wg.Done()
			auth, _ := ssh.ParseKeyFile(kfp)
			ssh.Tunnel = &ssh.TunnelRecipe{
				Auth:        auth,
				GatewayUser: cfg.Cluster.Owner,
				GatewayHost: gatewayStr,
				DialAddr:    remoteHost,
				BindAddr:    localHost,
				Log:         logger,
			}
			if err := ssh.Tunnel.Forward(ctx); err != nil {
				log.Printf("failed to forward tunnel - %s -> %s: %v", remoteHost, localHost, err)
				stop()
			}
		}(cfg.Cluster.KeyFile, gatewayAddr, chronosRemoteHost, endpoint.LocalHost.Chronos())
		wg.Wait()

		return nil
	},
}
