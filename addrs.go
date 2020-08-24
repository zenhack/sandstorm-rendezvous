package main

import (
	"fmt"
	"os"
)

func getenvDefault(v, def string) string {
	val := os.Getenv(v)
	if val == "" {
		return def
	}
	return val
}

func tmuxPath() string {
	return getenvDefault("TMUX_SOCKET", fmt.Sprintf("/tmp/tmux-%d/default", os.Getuid()))
}

func sandstormAddr() string {
	return getenvDefault("SANDSTORM_ADDR", ":6080")
}
