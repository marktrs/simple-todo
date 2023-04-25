package logger

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

var Log = log.Logger

func init() {
	// Persist logs to a file for later analysis and debugging if server was shuted down.
	runLogFile, _ := os.OpenFile(
		generateLogFilePath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)

	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
}

func HttpLogger(handler http.Handler) http.Handler {
	c := alice.New()
	c = c.Append(hlog.NewHandler(Log))
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

func generateLogFilePath() string {
	return filepath.Clean(strings.Join([]string{
		"./temp/",
		time.Now().Format("2006-01-02_15:04:05"),
		".log",
	}, ""))
}
