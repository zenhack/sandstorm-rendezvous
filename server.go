package main

import (
	"net/http"
)

func serverMain() {
	panic(http.ListenAndServe(":8000", NewWebServer()))
}
