package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var Logger *slog.Logger

func init() {
	flog, err := os.OpenFile(filepath.Join("log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("error opening file: %v", err))
	}
	Logger = slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, flog), nil))
}
