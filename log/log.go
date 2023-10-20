package log

import (
	"log/slog"
	"os"
)

var L = slog.New(
	slog.NewTextHandler(os.Stdout, nil))

// slog.NewJSONHandler(os.Stdout, nil))
