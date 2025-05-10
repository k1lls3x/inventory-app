package logs

import (

	"io"
	"log"
	"os"
	"path/filepath"
)

func Init() {
	logPath := filepath.Join("backend", "logs", "app.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("❌ лог-файл не открыть: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))

}
