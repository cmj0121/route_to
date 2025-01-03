package route_to

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/gin-gonic/gin"
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

func (r *RouteTo) Server() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()
	route.Any("/*path", r.serve)

	svc := &http.Server{
		Addr:    ":8080",
		Handler: route,
	}
	return svc
}

func (r *RouteTo) serve(c *gin.Context) {
	request := r.buildRequest(c)
	if request == nil {
		c.Status(http.StatusNotFound)
		return
	}

	client := &http.Client{}
	switch response, err := client.Do(request); err {
	case nil:
		defer response.Body.Close()

		for key, values := range response.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}
		r.postServe(c)

		if _, err := io.Copy(c.Writer, response.Body); err != nil {
			log.Error().Err(err).Msg("failed to copy the response")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(response.StatusCode)
	default:
		log.Error().Err(err).Msg("failed to send the request")
		c.Status(http.StatusInternalServerError)
	}
}

func (r *RouteTo) postServe(c *gin.Context) {
	if r.CORS {
		log.Debug().Msg("setting the CORS header")

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Max-Age", "86400")
	}
}

func (r *RouteTo) buildRequest(c *gin.Context) *http.Request {
	method := c.Request.Method
	endpoint := c.Param("path")

	switch endpoint {
	case "", "/":
		log.Debug().Msg("no endpoint specified")
		return nil
	default:
		endpoint = fmt.Sprintf("https://%s", endpoint[1:])
	}

	request, err := http.NewRequest(method, endpoint, c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to create the request")
		return nil
	}

	request.URL.RawQuery = c.Request.URL.RawQuery

	for key, values := range c.Request.Header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	log.Debug().Str("method", method).Str("endpoint", endpoint).Msg("building the request")
	return request
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
