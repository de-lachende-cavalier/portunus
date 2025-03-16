package main

import (
	"github.com/de-lachende-cavalier/portunus/cmd"
	"github.com/de-lachende-cavalier/portunus/pkg/logger"
)

func main() {
	// Initialize logger with pretty console output
	logger.Init("info", true)

	// Execute the root command
	cmd.Execute()
}
