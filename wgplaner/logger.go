package wgplaner

import (
	"os"

	"github.com/op/go-logging"
)

func init() {
	logFormat := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.999} [%{level:.4s}] %{module:5s} %{shortfunc:14s} ` +
			`▶ %{color:reset}%{message}`,
	)
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.NewBackendFormatter(logBackend, logFormat)

	logging.SetBackend(formatter)
}
