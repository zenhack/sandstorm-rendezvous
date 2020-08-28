package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/rpc"
)

func dialGrain(ctx context.Context, urlStr string, bootstrap *capnp.Client) *rpc.Conn {
	conn, _, err := (&websocket.Dialer{}).DialContext(ctx, urlStr, http.Header{})
	if err != nil {
		log.Fatalf("Connecting to grain: %v", err)
	}

	return rpc.NewConn(websocketTransport{conn}, &rpc.Options{
		BootstrapClient: bootstrap,
	})
}
