package main

import (
	"github.com/rs/zerolog/log"
	"github.com/cmj0121/route_to"
)

func main() {
	r := route_to.New()
	if err := r.ParseAndRun(); err != nil {
		log.Fatal().Err(err).Msg("failed to run the command")
	}
}
