package main

import (
	"context"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"github.com/gorilla/websocket"
)

func serveCapnp(ctx context.Context, wsConn *websocket.Conn, bootstrap capnp.Client) {
	transport := websocketTransport{wsConn}
	rpcConn := rpc.NewConn(transport, &rpc.Options{
		BootstrapClient: bootstrap,
	})
	select {
	case <-ctx.Done():
		rpcConn.Close()
	case <-rpcConn.Done():
		wsConn.Close()
	}
}
