package middleware

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/logger"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.StatusCode = statusCode
}

func HandleHTTPLogger(c *fiber.Ctx) error {
	startTime := time.Now()
	logger := logger.Log
	reqId := uuid.New().String()

	rq := c.Request()
	rs := c.Response()

	// Add request id to request header
	rq.Header.Set("RequestID", reqId)
	rq.Header.Set("Start", startTime.String())

	// Add request id to response header
	c.Set("Request-ID", reqId)

	logger.
		Info().
		Str("start", startTime.String()).
		Str("request_id", reqId).
		Str("protocol", "http").
		Str("method", string(rq.Header.Method())).
		Str("path", string(rq.RequestURI())).
		Str("status_text", http.StatusText(rs.StatusCode())).
		Msg("received a HTTP request")

	return c.Next()
}
