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
	conn, ln := dialGrain(ctx, urlStr)

	tmuxEndpoint := ip.TcpPort_ServerToClient(streamEndpoint{
		Network: "unix",
		Addr:    tmuxPath(),
	}, &server.Policy{})

	sandstormEndpoint := ip.TcpPort_ServerToClient(streamEndpoint{
		Network: "tcp",
		Addr:    sandstormAddr(),
	}, &server.Policy{})

	fut1, rel := ln.Bind(ctx, func(p LocalNetwork_bind_Params) error {
		info, err := p.NewInfo()
		if err != nil {
			return err
		}
		info.SetName("tmux")
		p.SetPort(tmuxEndpoint)
		return nil
	})

	fut2, rel := ln.Bind(ctx, func(p LocalNetwork_bind_Params) error {
		info, err := p.NewInfo()
		if err != nil {
			return err
		}
		info.SetName("sandstorm")
		p.SetPort(sandstormEndpoint)
		return nil
	})
	defer rel()
	h, err := fut1.Struct()
	if err != nil {
		log.Printf("Error binding tmux: ", err)
	} else {
		defer h.Handle().Client.Release()
	}
	h, err = fut2.Struct()
	if err != nil {
		log.Printf("Error binding sandstorm: ", err)
	} else {
		defer h.Handle().Client.Release()
	}
	log.Print("Listening!")
	<-conn.Done()
}

type streamEndpoint struct {
	Network, Addr string
}

func (ep streamEndpoint) Connect(ctx context.Context, p ip.TcpPort_connect) error {
	log.Println("Got connection")
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
