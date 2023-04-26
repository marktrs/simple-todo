package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

func generateLogFilePath() string {
	return filepath.Clean(strings.Join([]string{
		"./temp/",
		time.Now().Format("2006-01-02_15:04:05"),
		".log",
	}, ""))
}
