package main

import (
	"flag"
	"log"
	"os"
)

const usageMsg = `
CLI tool for message broacasting using websockets

Usage:

	broadcast-server <command> [arguments]

The commands are:

	start		start a websocket server
	connect		connect to a websocket server

Use "broadcast-server [command] -help|-h for more information about a command"

`

func main() {
	log.SetFlags(0)

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startPort := startCmd.Int("port", 8080, "listen for client connections on specified port")

	connectCmd := flag.NewFlagSet("connect", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Fatal(usageMsg)
	}

	switch os.Args[1] {
	case "start":
		if err := startCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		if err := StartServer(*startPort); err != nil {
			log.Fatalln("failed to start server")
		}
	case "connect":
		if err := connectCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		if err := ConnectClient(); err != nil {
			log.Fatal(err)
		}
	case "help":
		log.Fatal(usageMsg)
	default:
		log.Fatalf(
			"Unexpected subcommand \"%v\".\n\nUse \"broadcast-server help\" for more information about all subcommands.\n",
			os.Args[1],
		)
	}
}
