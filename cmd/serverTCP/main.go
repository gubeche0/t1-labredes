package main

import (
	"flag"
	"os"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/chat"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	MessagePort = flag.Int("message-port", 9000, "Port to recive message")
	CommandPort = flag.Int("command-port", 9001, "Port to recive command")
	User        = flag.String("user", "", "User to connect")

	Listen = flag.String("listen", "127.0.0.1", "Listen to connect")
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Parse()

	serv := chat.ServerTCP{
		Address:     *Listen,
		CommandPort: *CommandPort,
		MessagePort: *MessagePort,
	}

	serv.StartListenAndServer()
}
