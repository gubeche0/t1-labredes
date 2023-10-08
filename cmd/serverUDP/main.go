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
	Debug       = flag.Bool("debug", false, "Debug mode")

	Listen = flag.String("listen", "127.0.0.1", "Listen to connect")
)

func main() {
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode enabled")
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	serv := chat.ServerUDP{
		Address:     *Listen,
		CommandPort: *CommandPort,
		MessagePort: *MessagePort,
	}

	serv.StartListenAndServer()
}
