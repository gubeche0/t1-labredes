package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/gubeche0/raw-socket-t1-labredes/internal/chat"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	MessagePort = flag.Int("message-port", 9000, "Port to send message")
	CommandPort = flag.Int("command-port", 9001, "Port to send command")
	User        = flag.String("user", "", "User to connect")
	Debug       = flag.Bool("debug", false, "Debug mode")

	Destination = flag.String("destination", "localhost", "Destination to send message")
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

	if *User == "" {
		// Generate random user
		*User = fmt.Sprintf("Anonymous_%d", rand.Intn(1000))
	}

	client := chat.ClientChat{
		User:        *User,
		Address:     *Destination,
		MessagePort: *MessagePort,
		CommandPort: *CommandPort,
	}

	client.Connect()
}
