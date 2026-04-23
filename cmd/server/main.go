package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codebase-brain/internal/config"
	"github.com/codebase-brain/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.Default()
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		if loaded, err := config.Load(configPath); err == nil {
			cfg = loaded
		}
	}
	cfg = config.FromEnv()

	s := server.NewMCPServer(
		"AI Coder Context Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
	)

	mcp.RegisterTools(s, cfg)

	log.Printf("Starting AI Coder Context Server on port %d", cfg.Port)
	log.Printf("Project path: %s", cfg.ProjectPath)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}