package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"zombiezen.com/go/capnproto2/rpc"
)

func dialGrain(ctx context.Context, urlStr string) (*rpc.Conn, LocalNetwork) {
	conn, _, err := (&websocket.Dialer{}).DialContext(ctx, urlStr, http.Header{})
	if err != nil {
		log.Fatal("Connecting to grain: %v", err)
	}

	rpcConn := rpc.NewConn(websocketTransport{conn}, nil)
	return rpcConn, LocalNetwork{rpcConn.Bootstrap(ctx)}
}
