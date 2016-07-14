package rpclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateAccuont(t *testing.T) {

}

// CaseHandler represents one test case, which will be invoked by MockServer.
type CaseHandler func() (*httptest.ResponseRecorder, *http.Request)

// MockServer mock server state for various test cases,
// people can fake the server's state by fullfill the Server interface, and
// define various request cases by defining functions that match the signature of
// CaseHandler.
func MockServer(cli Client, fs CaseHandler) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := NewRouter(cli)
	w, r := fs()
	router.ServeHTTP(w, r)
	return w
}

// HttpRequestCase is used to create REST api test cases.
func HttpRequestCase(method string, url string, body io.Reader) CaseHandler {
	return func() (*httptest.ResponseRecorder, *http.Request) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(method, url, body)
		if err != nil {
			panic(err)
		}
		switch method {
		case "POST":
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		return w, r
	}
}
