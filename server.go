package route_to

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// The interface for the HTTP client and used to make the requests and make
// the function testable.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

var (
	default_client HTTPClient = &http.Client{}
)

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

	switch response, err := default_client.Do(request); err {
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
