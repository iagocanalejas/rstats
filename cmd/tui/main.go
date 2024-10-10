package main

import (
	"log"
	"os"

	"github.com/iagocanalejas/rstats/internal/utils/assert"
	"github.com/iagocanalejas/rstats/pkg/tui"
)

func main() {
	setupFileLogger()
	if err := tui.BuildApp().App.Run(); err != nil {
		panic(err)
	}
}

func setupFileLogger() {
	// create a log file as with the tui we can't see the logs
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	assert.NoError(err, "opening log file")

	log.SetOutput(f)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
