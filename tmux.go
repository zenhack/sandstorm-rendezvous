package main

import (
	"fmt"
	"os"
)

func tmuxPath() string {
	return fmt.Sprintf("/tmp/tmux-%d/default", os.Getuid())
}
