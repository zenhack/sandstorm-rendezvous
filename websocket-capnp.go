package main

import (
	"context"
	"log"

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
	msg, _, err := capnp.NewMessage(arena)
	if err != nil {
		return rpccp.Message{}, func() error { return nil }, func() {}, err
	}
	send = func() error {
		data, err := msg.Marshal()
		if err != nil {
			log.Print("Error getting segment: ", err)
			return err
		}
		log.Print("About to send message of size: ", len(data))

		return t.conn.WriteMessage(websocket.BinaryMessage, data)
	}
	release := func() {}
	seg, _ := msg.Segment(0)
	rpcMsg, err := rpccp.NewRootMessage(seg)
	return rpcMsg, send, release, err
}

func (t websocketTransport) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	log.Println("RecvMessage()")
	var (
		typ  int
		data []byte
		err  error
	)
	for ctx.Err() == nil && typ != websocket.BinaryMessage {
		typ, data, err = t.conn.ReadMessage()
		log.Printf("typ = %v, len(data) = %v, err = %v", typ, len(data), err)
		if err != nil {
			return rpccp.Message{}, func() {}, err
		}
		if typ == websocket.PingMessage {
			t.conn.WriteMessage(websocket.PongMessage, nil)
		}
	}
	if err = ctx.Err(); err != nil {
		return rpccp.Message{}, func() {}, err
	}
	log.Println("Got message of size %v", len(data))

	msg, err := capnp.Unmarshal(data)
	if err != nil {
		return rpccp.Message{}, func() {}, err
	}
	rpcMsg, err := rpccp.ReadRootMessage(msg)
	return rpcMsg, func() {}, err
}

func (t websocketTransport) Close() error {
	return t.conn.Close()
}
