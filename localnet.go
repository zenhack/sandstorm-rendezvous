package main

import (
	"context"
	"errors"
	"log"

	"zenhack.net/go/sandstorm/capnp/ip"
	"zenhack.net/go/sandstorm/capnp/util"
	"zombiezen.com/go/capnproto2/server"
)

var _ LocalNetwork_Server = &LocalNetworkImpl{}

var (
	ErrAlreadyBound = errors.New("The name is already bound")
	ErrNotFound     = errors.New("No such port.")
)

type LocalNetworkImpl struct {
	mu    chan struct{}
	ports map[string]ip.TcpPort
}

func newLocalNetwork() *LocalNetworkImpl {
	ret := &LocalNetworkImpl{
		mu:    make(chan struct{}, 1),
		ports: make(map[string]ip.TcpPort),
	}
	ret.unlock()
	return ret
}

func (ln *LocalNetworkImpl) unlock() {
	log.Println("Unlocking network")
	ln.mu <- struct{}{}
}

func (ln *LocalNetworkImpl) lock() {
	log.Println("Locking network")
	<-ln.mu
}

func (ln *LocalNetworkImpl) Bind(ctx context.Context, p LocalNetwork_bind) error {
	log.Print("Bind()")
	params := p.Args()
	port := params.Port()

	info, err := params.Info()
	if err != nil {
		return err
	}

	name, err := info.Name()
	if err != nil {
		return err
	}
	log.Print("name = ", name)

	res, err := p.AllocResults()
	if err != nil {
		return err
	}

	ln.lock()
	defer ln.unlock()

	if _, ok := ln.ports[name]; ok {
		return ErrAlreadyBound
	}
	ln.ports[name] = port
	port.Client.AddRef()

	handle := util.Handle_ServerToClient(&portDropHandle{
		name: name,
		ln:   ln,
	}, &server.Policy{})
	res.SetHandle(handle)
	log.Print("ok.")
	return nil
}

func (ln *LocalNetworkImpl) Resolve(ctx context.Context, p LocalNetwork_resolve) error {
	log.Print("Resolve()")
	res, err := p.AllocResults()
	if err != nil {
		return err
	}
	name, err := p.Args().Name()
	if err != nil {
		return err
	}
	log.Print("name = ", name)

	ln.lock()
	defer ln.unlock()
	port, ok := ln.ports[name]
	if !ok {
		return ErrNotFound
	}
	res.SetPort(port)
	log.Print("ok.")
	return nil
}

type portDropHandle struct {
	name string
	ln   *LocalNetworkImpl
}

func (*portDropHandle) Ping(context.Context, util.Handle_ping) error {
	return nil
}

func (h *portDropHandle) Shutdown() error {
	h.ln.lock()
	defer h.ln.unlock()
	log.Print("Dropping name = ", h.name)
	port := h.ln.ports[h.name]
	delete(h.ln.ports, h.name)
	port.Client.Release()
	return nil
}
