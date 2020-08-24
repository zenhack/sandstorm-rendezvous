module zenhack.net/go/sandstorm-rendezvous

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	zenhack.net/go/sandstorm v0.0.0-20200807223653-d169734aeb58
	zombiezen.com/go/capnproto2 v2.18.0+incompatible
)

replace zombiezen.com/go/capnproto2 => /home/isd/src/foreign/go-capnproto2
