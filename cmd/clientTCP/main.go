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

	Destination = flag.String("destination", "localhost", "Destination to send message")
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Parse()
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
