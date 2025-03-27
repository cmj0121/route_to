package route_to

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	callCount int
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.callCount++

	resp := httptest.NewRecorder()
	resp.WriteHeader(200)
	_, _ = resp.Write([]byte("Hello, World!"))

	return resp.Result(), nil
}

func (m *MockHTTPClient) Reset() {
	m.callCount = 0
}

func TestServer(t *testing.T) {
	c := &MockHTTPClient{}
	default_client = c

	cases := []string{
		"/example.com",
		"/example.com/",
		"/example.com/path/to/resource",
	}

	for _, link := range cases {
		c.Reset()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", link, nil)

		svc := New().Server()
		svc.Handler.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, c.callCount)
	}
}
