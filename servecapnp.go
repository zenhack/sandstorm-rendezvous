package main

import (
	"context"

	"github.com/gorilla/websocket"
	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/rpc"
)

func serveCapnp(ctx context.Context, wsConn *websocket.Conn, bootstrap *capnp.Client) {
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
