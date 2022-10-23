package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"

	"zenhack.net/go/sandstorm/capnp/ip"
	"zenhack.net/go/sandstorm/capnp/util"
)

func NewWebServer() http.Handler {
	hostClient := ip.TcpPort{}
	hostLock := &sync.Mutex{}

	up := &websocket.Upgrader{}
	r := http.NewServeMux()
	r.HandleFunc("/host.socket", func(w http.ResponseWriter, req *http.Request) {
		wsConn, err := up.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Error upgrading websocket:", err)
			return
		}
		func() {
			hostLock.Lock()
			defer hostLock.Unlock()
			if (hostClient != ip.TcpPort{}) {
				log.Println("Host client already connected; rejecting.")
				w.WriteHeader(500)
				return
			}
			rpcConn := rpc.NewConn(websocketTransport{wsConn}, nil)
			client := rpcConn.Bootstrap(req.Context())
			hostClient = ip.TcpPort(client)
			go func() {
				<-req.Context().Done()
				hostLock.Lock()
				defer hostLock.Unlock()
				if capnp.Client(hostClient) == client {
					hostClient = ip.TcpPort{}
				}
			}()
		}()
		<-req.Context().Done()
	})
	r.HandleFunc("/guest.socket", func(w http.ResponseWriter, req *http.Request) {
		conn, err := up.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Error upgrading websocket:", err)
			return
		}
		defer conn.Close()
		serveGuest(req.Context(), conn, hostClient)
	})
	r.Handle("/", http.FileServer(http.Dir("static")))
	r.HandleFunc("/sandstorm-rendezvous", func(w http.ResponseWriter, req *http.Request) {
		f, err := os.Open("/sandstorm-rendezvous")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		defer f.Close()
		io.Copy(w, f)
	})
	return r
}

func serveGuest(ctx context.Context, conn *websocket.Conn, port ip.TcpPort) {
	res, release := port.Connect(ctx, func(p ip.TcpPort_connect_Params) error {
		p.SetDownstream(util.ByteStream_ServerToClient(
			websocketByteStream{conn},
		))
		return nil
	})
	defer release()
	upstream := res.Upstream()
	ctx, cancel := context.WithCancel(ctx)
	errCh := make(chan func() error, 10) // buffer size is arbitrary.
	go func() {
		defer cancel()
		for {
			select {
			case errFn := <-errCh:
				err := errFn()
				if err != nil {
					log.Println("Error writing to bytestream:", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	for ctx.Err() == nil {
		typ, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("reading from websocket: ", err)
			return
		}
		switch typ {
		case websocket.CloseMessage:
			res, release := upstream.Done(ctx, func(util.ByteStream_done_Params) error {
				return nil
			})
			release()
			_, err = res.Struct()
			if err != nil {
				log.Println("ByteStream.done():", err)
			}
			return
		case websocket.BinaryMessage:
			res, release := upstream.Write(ctx, func(p util.ByteStream_write_Params) error {
				p.SetData(data)
				return nil
			})
			errCh <- func() error {
				_, err := res.Struct()
				release()
				return err
			}
		}
	}
}

type websocketByteStream struct {
	conn *websocket.Conn
}

func (s websocketByteStream) Write(ctx context.Context, p util.ByteStream_write) error {
	data, err := p.Args().Data()
	if err != nil {
		return err
	}
	return s.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (s websocketByteStream) Done(context.Context, util.ByteStream_done) error {
	return s.conn.Close()
}

func (s websocketByteStream) ExpectSize(context.Context, util.ByteStream_expectSize) error {
	return nil
}
