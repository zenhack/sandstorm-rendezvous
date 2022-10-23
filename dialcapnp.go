package main

import (
	"context"
	"log"
	"net/http"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"github.com/gorilla/websocket"
)

func dialGrain(ctx context.Context, urlStr string, bootstrap capnp.Client) *rpc.Conn {
	conn, _, err := (&websocket.Dialer{}).DialContext(ctx, urlStr, http.Header{})
	if err != nil {
		log.Fatalf("Connecting to grain: %v", err)
	}

	return rpc.NewConn(websocketTransport{conn}, &rpc.Options{
		BootstrapClient: bootstrap,
	})
}
