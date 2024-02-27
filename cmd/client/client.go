package main

import (
	"bufio"
	"fmt"
	"github.com/stas9132/GophKeeper/internal/client"
	"os"
	"strings"
)

func shell() {
	cl := client.NewClient()
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		flds := strings.Fields(s.Text())
		if len(flds) == 0 {
			continue
		}
		cmd := flds[0]
		switch cmd {
		case "exit":
			return
		case "dial":
			if err := cl.Dial(); err != nil {
				fmt.Println(err)
				continue
			}
		case "close":
			if err := cl.Close(); err != nil {
				fmt.Println(err)
				continue
			}
		case "health":
			if err := cl.Health(); err != nil {
				fmt.Println(err)
				continue
			}
		default:
			fmt.Println("unknown command")
			continue
		}
		fmt.Println("command complete successfully")
	}
}

func main() {
	shell()
}
