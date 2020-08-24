package main

import (
	"net/http"
)

func serverMain() {
	ln := LocalNetwork_ServerToClient(newLocalNetwork(), nil)
	webSrv := NewWebServer(ln)
	panic(http.ListenAndServe(":8000", webSrv))
}
