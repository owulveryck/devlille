package main

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleDemoContent(t *testing.T) {
	ctx := context.Background()
	request := mcp.ReadResourceRequest{}

	contents, err := handleDemoContent(ctx, request)
	if err != nil {
		t.Fatalf("handleDemoContent failed: %v", err)
	}

	if len(contents) != 1 {
		t.Fatalf("expected 1 content, got %d", len(contents))
	}

	content := contents[0]
	
	// Check that text content contains concatenated text from the database
	textContent, ok := content.(mcp.TextResourceContents)
	if !ok {
		t.Fatal("expected TextResourceContents")
	}

	if textContent.URI != "demo://content" {
		t.Errorf("expected URI 'demo://content', got '%s'", textContent.URI)
	}

	if textContent.MIMEType != "text/plain" {
		t.Errorf("expected MIME type 'text/plain', got '%s'", textContent.MIMEType)
	}

	expectedText := "aha je veux comprendre le protocol MCP dsadas dsdas"
	if textContent.Text != expectedText {
		t.Errorf("expected text '%s', got '%s'", expectedText, textContent.Text)
	}
}