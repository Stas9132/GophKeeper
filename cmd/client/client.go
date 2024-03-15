package main

import (
	"bufio"
	"fmt"
	"github.com/stas9132/GophKeeper/internal/client"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"log"
	"os"
	"strings"
)

func shell(l logger.Logger) {
	cl, err := client.NewClient(l)
	if err != nil {
		log.Fatalln(err)
	}
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
		case "health":
			if err = cl.Health(); err != nil {
				fmt.Println(err)
				continue
			}
		case "register":
			if err = cl.Register(flds); err != nil {
				fmt.Println(err)
				continue
			}
		case "login":
			if err = cl.Login(flds); err != nil {
				fmt.Println(err)
				continue
			}
		case "logout":
			if err = cl.Logout(); err != nil {
				fmt.Println(err)
				continue
			}
		case "put":
			if err = cl.Put(flds); err != nil {
				fmt.Println(err)
				continue
			}
		case "get":
			data, err := cl.Get(flds)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(data)
		case "list":
			keys, err := cl.List()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for i, key := range keys {
				fmt.Println(i, key)
			}
		case "sync":
			if err = cl.Sync(); err != nil {
				fmt.Println(err)
				continue
			}
		case "help":
			fmt.Println("Valid commands:\nregister\nlogin\nlogout\nhealth\nput\nget\nsync\nexit")
		default:
			fmt.Println("unknown command")
			continue
		}
		fmt.Println("command complete successfully")
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.Init()

	if config.PrintVersion {
		log.Println("Version:", config.Version)
		log.Println("Build date:", config.BuildDate)
		return
	}

	l := logger.NewSlogLogger()
	shell(l)
}
