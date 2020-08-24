package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"zenhack.net/go/sandstorm/capnp/ip"
	"zenhack.net/go/sandstorm/exp/util/bytestream"
)

func connectMain(urlStr string) {
	ctx := context.Background()
	conn, ln := dialGrain(ctx, urlStr)
	go runProxy(ctx, ln, "tmux", "unix", tmuxPath())
	go runProxy(ctx, ln, "sandstorm", "tcp", ":6080")
	select {
	case <-conn.Done():
	case <-ctx.Done():
	}
}

func runProxy(ctx context.Context, ln LocalNetwork, name, network, addr string) {
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalf("Listening for %q: %v", name, err)
	}
	res, release := ln.Resolve(ctx, func(p LocalNetwork_resolve_Params) error {
		p.SetName(name)
		return nil
	})
	defer release()
	port := res.Port()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error: Accept(): %v", err)
			time.Sleep(time.Second)
			continue
		}

		go handleConn(ctx, conn, port)
	}
}

func handleConn(ctx context.Context, conn net.Conn, port ip.TcpPort) {
	defer conn.Close()
	res, release := port.Connect(ctx, func(p ip.TcpPort_connect_Params) error {
		p.SetDownstream(bytestream.FromWriteCloser(conn, nil))
		return nil
	})
	defer release()
	w := bytestream.ToWriteCloser(ctx, res.Upstream())
	io.Copy(w, conn)
}
