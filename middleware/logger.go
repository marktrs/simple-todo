package middleware

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/marktrs/simple-todo/logger"
	"github.com/rs/zerolog/hlog"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func HTTPLogger(handler http.Handler) http.Handler {
	c := alice.New()
	c = c.Append(hlog.NewHandler(logger.Log))
	c = c.Append(hlog.RequestIDHandler("request_id", "Request-Id"))
	h := c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, r)
		duration := time.Since(startTime)

		logger := hlog.FromRequest(r).Info()
		if rec.StatusCode != http.StatusOK {
			logger = hlog.FromRequest(r).Error().Bytes("body", rec.Body)
		}

		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	}))

	return h
}
