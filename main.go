package main

import (
	"CloudPhoto/cmd/daemon"
	"CloudPhoto/cmd/server"
)

func main() {
	server.Start()
	daemon.Daemon()
	select {}
}
