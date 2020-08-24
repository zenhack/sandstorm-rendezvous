package main

import (
	"context"

	"github.com/gorilla/websocket"
	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/rpc"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

type websocketTransport struct {
	conn *websocket.Conn
}

var _ rpc.Transport = websocketTransport{}

func (t websocketTransport) NewMessage(ctx context.Context) (
	_ rpccp.Message,
	send func() error,
	_ capnp.ReleaseFunc,
	_ error,
) {
	arena := capnp.SingleSegment(nil)
	msg := &capnp.Message{Arena: arena}
	send = func() error {
		// SingleSegment never returns an error when accessing segment 0:
		data, _ := arena.Data(0)

		return t.conn.WriteMessage(websocket.BinaryMessage, data)
	}
	release := func() {}
	seg, _ := msg.Segment(0)
	rpcMsg, err := rpccp.NewRootMessage(seg)
	return rpcMsg, send, release, err
}

func (t websocketTransport) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	var (
		typ  int
		data []byte
		err  error
	)
	for ctx.Err() == nil {
		typ, data, err = t.conn.ReadMessage()
		if err != nil {
			return rpccp.Message{}, func() {}, err
		}
		switch typ {
		case websocket.PingMessage:
			t.conn.WriteMessage(websocket.PongMessage, nil)
		case websocket.BinaryMessage:
			break
		default:
			continue
		}
	}
	if err = ctx.Err(); err != nil {
		return rpccp.Message{}, func() {}, err
	}

	msg := &capnp.Message{Arena: capnp.SingleSegment(data)}
	rpcMsg, err := rpccp.ReadRootMessage(msg)
	return rpcMsg, func() {}, err
}

func (t websocketTransport) Close() error {
	return t.conn.Close()
}
