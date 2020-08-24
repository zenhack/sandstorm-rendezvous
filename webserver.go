package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var indexHtml = []byte(`
<!doctype html>
<html>
	<head>
		<meta charset="utf-8" />
		<title>Rendezvous</title>
		<script src="/offer-iframe.js"></script>
	</head>
	<body>
		<iframe style="width: 100%; height: 55px; marign: 0; border: 0;" id="offer-iframe"></iframe>
	</body>
</html>
`)

var offerIframeJS = []byte(`
window.addEventListener('message', function(event) {
	if (event.data.rpcId !== "0") {
		return;
	}
	if (event.data.error) {
		console.log("ERROR: " + event.data.error);
	} else {
		const el = document.getElementById("offer-iframe");
		el.setAttribute("src", event.data.uri);
	}
});
document.addEventListener('DOMContentLoaded', function() {
	const template = window.location.protocol.replace('http', 'ws') +
		"//$API_HOST/.sandstorm-token/$API_TOKEN/socket";
	window.parent.postMessage({renderTemplate: {
		rpcId: "0",
		template: template,
		clipboardButton: 'left'
	}}, "*");
})
`)

func NewWebServer(ln LocalNetwork) http.Handler {
	up := &websocket.Upgrader{}
	r := mux.NewRouter()
	r.HandleFunc("/socket", func(w http.ResponseWriter, req *http.Request) {
		conn, err := up.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Error upgrading websocket:", err)
			w.WriteHeader(500)
			return
		}
		serveCapnp(req.Context(), conn, ln.Client)
	})
	r.HandleFunc("/offer-iframe.js", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(offerIframeJS)
	})
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHtml)
	})
	return r
}
