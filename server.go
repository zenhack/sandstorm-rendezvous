package main

import (
	"net/http"
	"zombiezen.com/go/capnproto2/server"
)

func serverMain() {
	ln := LocalNetwork_ServerToClient(newLocalNetwork(), &server.Policy{})
	webSrv := NewWebServer(ln)
	panic(http.ListenAndServe(":8000", webSrv))
}
