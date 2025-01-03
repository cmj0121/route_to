package main

import (
	"github.com/cmj0121/route_to"
	"github.com/rs/zerolog/log"
)

func main() {
	r := route_to.New()
	if err := r.ParseAndRun(); err != nil {
		log.Fatal().Err(err).Msg("failed to run the command")
	}
}
