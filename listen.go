package main

import (
	"context"
	"io"
	"log"
	"net"

	"zenhack.net/go/sandstorm/capnp/ip"
	"zenhack.net/go/sandstorm/capnp/util"
	"zenhack.net/go/sandstorm/exp/util/bytestream"
	"zombiezen.com/go/capnproto2/server"
)

func listenMain(urlStr string) {
	ctx := context.Background()
	vncEndpoint := ip.TcpPort_ServerToClient(streamEndpoint{
		Network: "tcp",
		Addr:    "127.0.0.1:5901",
	}, &server.Policy{})

	conn := dialGrain(ctx, urlStr, vncEndpoint.Client)

	log.Print("Listening...")
	<-conn.Done()
}

type streamEndpoint struct {
	Network, Addr string
}

func (ep streamEndpoint) Connect(ctx context.Context, p ip.TcpPort_connect) error {
	log.Println("Got connection.")
	res, err := p.AllocResults()
	if err != nil {
		return err
	}

	conn, err := net.Dial(ep.Network, ep.Addr)
	if err != nil {
		return err
	}

	res.SetUpstream(bytestream.FromWriteCloser(conn, &server.Policy{}))
	downstream := p.Args().Downstream()
	downstream = util.ByteStream{downstream.Client.AddRef()}

	w := bytestream.ToWriteCloser(context.TODO(), downstream)
	go func() {
		defer conn.Close()
		defer downstream.Client.Release()
		io.Copy(w, conn)
	}()
	return nil
}
