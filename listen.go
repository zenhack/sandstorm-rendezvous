package main

import (
	"context"
	"io"
	"net"

	"zenhack.net/go/sandstorm/capnp/ip"
	"zenhack.net/go/sandstorm/exp/util/bytestream"
)

func listenMain(urlStr string) {
	ctx := context.Background()
	conn, ln := dialGrain(ctx, urlStr)

	tmuxEndpoint := ip.TcpPort_ServerToClient(streamEndpoint{
		Network: "unix",
		Addr:    tmuxPath(),
	}, nil)

	sandstormEndpoint := ip.TcpPort_ServerToClient(streamEndpoint{
		Network: "tcp",
		Addr:    ":6080",
	}, nil)

	ln.Bind(ctx, func(p LocalNetwork_bind_Params) error {
		info, err := p.NewInfo()
		if err != nil {
			return err
		}
		info.SetName("tmux")
		p.SetPort(tmuxEndpoint)
		return nil
	})

	ln.Bind(ctx, func(p LocalNetwork_bind_Params) error {
		info, err := p.NewInfo()
		if err != nil {
			return err
		}
		info.SetName("sandstorm")
		p.SetPort(sandstormEndpoint)
		return nil
	})

	<-conn.Done()
}

type streamEndpoint struct {
	Network, Addr string
}

func (ep streamEndpoint) Connect(ctx context.Context, p ip.TcpPort_connect) error {
	res, err := p.AllocResults()
	if err != nil {
		return err
	}

	conn, err := net.Dial(ep.Network, ep.Addr)
	if err != nil {
		return err
	}

	res.SetUpstream(bytestream.FromWriteCloser(conn, nil))
	downstream := p.Args().Downstream()
	downstream.Client.AddRef()

	w := bytestream.ToWriteCloser(context.TODO(), downstream)
	go func() {
		defer conn.Close()
		defer downstream.Client.Release()
		io.Copy(w, conn)
	}()
	return nil
}
