package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/tag-service/test"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

//tag-service | [GIN] 2018/02/15 - 16:48:52 |200 |        40.9µs | X-Forwarded-For: 172.20.0.1| Referer: Banish_The_Go_Path_man | User-Agent: HTTPie/0.9.9 | GET     /health

const (
	GoPathRevolt  = "Banish_The_Go_Path_man"
	XforwardedFor = "X-Forwarded-For"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestTestLogger(t *testing.T) {

	t.Logf("Given an Http request is performed to the below endpionts")
	{
		buffer := new(bytes.Buffer)
		router := gin.New()
		router.Use(LogWithWriter(buffer, "hello"))

		// all possible routes
		router.GET("/example", func(c *gin.Context) {})
		router.GET("/hello", func(c *gin.Context) {})
		router.POST("/example", func(c *gin.Context) {})
		router.PUT("/example", func(c *gin.Context) {})
		router.DELETE("/example", func(c *gin.Context) {})
		router.PATCH("/example", func(c *gin.Context) {})
		router.HEAD("/example", func(c *gin.Context) {})
		router.OPTIONS("/example", func(c *gin.Context) {})

		performRequest(router, "GET", "/example?a=100")

		if strings.Contains(buffer.String(), "X-Forwarded-For:       localhost") {
			t.Logf("\t\tThe X-Forwarded-For should have been captured in the log \"%s\". %v", buffer.String(), test.CheckMark)
		} else {
			t.Errorf("\t\tThe X-Forwarded-For should have been captured in the log \"%s\". %v", buffer.String(), test.BallotX)
		}

		if strings.Contains(buffer.String(), "Referer: Banish_The_Go_Path_man") {
			t.Logf("\t\tThe Referer should have been captured in the log: \"%s\". %v", buffer.String(), test.CheckMark)
		} else {
			t.Errorf("\t\tThe Referer should have been captured in the log:  \"%s\". %v", buffer.String(), test.BallotX)
		}

		if strings.Contains(buffer.String(), "User-Agent: Banish_The_Go_Path_man ") {
			t.Logf("\t\tThe User-Agent should have been captured in the log: \"%s\". %v", buffer.String(), test.CheckMark)
		} else {
			t.Errorf("\t\tThe User-Agent should have been captured in the log: \"%s\". %v", buffer.String(), test.BallotX)
		}

	}
}

func TestLogWithWriter(t *testing.T) {
	handlerFunc := Logger()
	if reflect.TypeOf(handlerFunc).Name() == "HandlerFunc" {
		t.Logf("\t\tThe logger middleware should have been HandlerFunc. %v", test.CheckMark)
	} else {
		t.Errorf("\t\tThe logger middleware should have been HandlerFunc. %v", test.BallotX)
	}
}

// Performs http request
func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Add(Referer, GoPathRevolt)
	req.Header.Add(UserAgent, GoPathRevolt)
	req.Header.Add(XforwardedFor, "localhost")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
