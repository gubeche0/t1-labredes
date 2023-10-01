package chat

import (
	"io"

	"github.com/rs/zerolog/log"
)

func checkErrorMessage(err error) bool {
	if err == io.EOF {
		log.Info().Msg("Connection closed")
		return false
	} else if err != nil {
		log.Err(err).Msg("Error to read message")
		return false
	}

	return true
}
