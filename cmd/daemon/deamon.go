package daemon

import (
	"CloudPhoto/cmd/server"
	"CloudPhoto/config"
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commands = map[string]func(){
	"stop":          server.Stop,
	"exist":         server.Stop,
	"reload-config": config.Read,
}

func Daemon() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("CloudPhoto>>>")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if command, ok := commands[cmd]; ok {
			command()
		} else {
			fmt.Println("command not found")
		}
	}
}
