package route_to

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The main instance of the RouteTo struct and used to launch the service.
type RouteTo struct {
	Verbose int `short:"v" type:"counter" help:"Set the verbose level of the command."`

	CORS bool `help:"Enable the CORS for the service."`
}

// Create a new instance of the RouteTo struct with the default values.
func New() *RouteTo {
	return &RouteTo{}
}

// Parse the command line arguments and run the service.
func (r *RouteTo) ParseAndRun() error {
	ctx := kong.Parse(r)
	return r.Run(ctx)
}

// Run the service with the given context.
func (r *RouteTo) Run(ctx *kong.Context) error {
	r.prologue()
	defer r.epilogue()

	return r.run()
}

func (r *RouteTo) run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	svc := r.Server()
	go func() {
		log.Info().Str("addr", svc.Addr).Msg("starting the service ...")

		if err := svc.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start the service")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutting down the service ...")
	return svc.Shutdown(ctx)
}
// setup everything before running the command
func (r *RouteTo) prologue() {
	// setup the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// set the verbose level
	switch r.Verbose {
	case 0:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	log.Info().Msg("starting the relink ...")
}

// cleanup everything after running the command
func (r *RouteTo) epilogue() {
	log.Info().Msg("finished the relink ...")
}
