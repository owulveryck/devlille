package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var dbPath = flag.String("db", "../DB/db.json", "Path to the database file")

type DBEntry struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	UserAgent string `json:"user_agent"`
	Timestamp string `json:"timestamp"`
}

func NewDemoMCPServer() *server.MCPServer {
	mcpServer := server.NewMCPServer(
		"demo-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithInstructions(`Ceci est un server MCP (Model Context Protocol) de demonstration dans le cadre de la conférence du DevLille.
Ce serveur expose les attentes des participants via l'URI demo://content ainsi qu'un prompt pour jouer la demo.`),
	)

	mcpServer.AddResource(
		mcp.NewResource("demo://content",
			"DemoContentResource",
			mcp.WithMIMEType("text/plain"),
		),
		handleDemoContent,
	)
	mcpServer.AddPrompt(GetPromptSummary(), GetPromptSummaryHandler)
	mcpServer.AddPrompt(GetPrompt(), GetPromptHandler)

	return mcpServer
}

func handleDemoContent(
	ctx context.Context,
	request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	// Read the DB file
	data, err := os.ReadFile(*dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database file: %w", err)
	}

	// Parse JSON
	var entries []DBEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Concatenate all text content
	var textContent strings.Builder
	for i, entry := range entries {
		if i > 0 {
			textContent.WriteString(" ")
		}
		textContent.WriteString(strings.TrimSpace(entry.Text))
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "demo://content",
			MIMEType: "text/plain",
			Text:     textContent.String()[:255],
		},
	}, nil
}

func main() {
	flag.Parse()

	mcpServer := NewDemoMCPServer()

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
