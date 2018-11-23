package main

import (
	"flag"
	"fmt"
	"github.com/eagledb/eagledb/config"
	"github.com/eagledb/eagledb/server"
	"os"
)

var confFile string

func init() {
	flag.StringVar(&confFile, "config", "/etc/eagledb/eagledb.toml", "eagledb configuration file")
	flag.Parse()
}

func main() {
	err := config.LoadFile(confFile)
	if err != nil {
		fmt.Printf("failed to load config file \"%s\": %s\n", confFile, err)
		os.Exit(1)
	}

	server := server.NewServer()

	err = server.Start()
	if err != nil {
		fmt.Println("failed to start eagledb:", err)
		os.Exit(1)
	}
}
