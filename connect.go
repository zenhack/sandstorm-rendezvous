package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"zenhack.net/go/sandstorm/capnp/ip"
	"zenhack.net/go/sandstorm/exp/util/bytestream"
	"zombiezen.com/go/capnproto2/server"
)

func connectMain(urlStr string) {
	ctx := context.Background()
	lnc1 := make(chan LocalNetwork, 1)
	lnc2 := make(chan LocalNetwork, 1)
	err := os.MkdirAll(fmt.Sprintf("/tmp/tmux-%d", os.Getuid()), 0750)
	if err != nil {
		log.Fatal("could not create socket dir: ", err)
	}
	go runProxy(ctx, lnc1, "tmux", "unix", tmuxPath())
	go runProxy(ctx, lnc2, "sandstorm", "tcp", sandstormAddr())
	conn, ln := dialGrain(ctx, urlStr)
	lnc1 <- ln
	lnc2 <- ln
	select {
	case <-conn.Done():
	case <-ctx.Done():
	}
}

func runProxy(ctx context.Context, lnc chan LocalNetwork, name, network, addr string) {
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalf("Listening for %q: %v", name, err)
	}
	if network == "unix" {
		err = os.Chmod(addr, 0660)
		if err != nil {
			log.Fatal("Setting socket permissions: ", err)
		}
	}
	defer l.Close()
	ln := <-lnc
	res, release := ln.Resolve(ctx, func(p LocalNetwork_resolve_Params) error {
		p.SetName(name)
		return nil
	})
	defer release()
	r, err := res.Struct()
	if err != nil {
		log.Fatal("Error: ln.Resolve(): ", err)
	}
	port := r.Port()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error: Accept(): %v", err)
			time.Sleep(time.Second)
			continue
		}
		log.Print("Got connection for ", name)

		go handleConn(ctx, conn, port)
	}
}

func handleConn(ctx context.Context, conn net.Conn, port ip.TcpPort) {
	defer conn.Close()
	res, release := port.Connect(ctx, func(p ip.TcpPort_connect_Params) error {
		p.SetDownstream(bytestream.FromWriteCloser(conn, &server.Policy{}))
		return nil
	})
	defer release()
	w := bytestream.ToWriteCloser(ctx, res.Upstream())
	if _, err := res.Struct(); err != nil {
		log.Print("TcpPort.Connect(): ", err)
		return
	}
	io.Copy(w, conn)
}
