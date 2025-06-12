package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPromptSummary() mcp.Prompt {
	return mcp.NewPrompt("getSummaryOfRequests",
		mcp.WithPromptDescription("Ce prompt extrait les attentes des participants de manière formattée."),
		mcp.WithArgument("input",
			mcp.ArgumentDescription("Les entrées brutes des attentes des participants"),
			mcp.RequiredArgument(),
		),
	)
}

func GetPromptSummaryHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	arguments := request.Params.Arguments
	return &mcp.GetPromptResult{
		Description: "Resultat de l'analyse",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: `Tu vas générer un rapport structuré des attentes des particpants.`,
				},
			},
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: arguments["input"],
				},
			},
		},
	}, nil
}

func GetPrompt() mcp.Prompt {
	return mcp.NewPrompt("workflow",
		mcp.WithPromptDescription("Ce prompt permet de lancer la démo."),
	)
}

func GetPromptHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Ce prompt permet d'éxecuter la demo",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: `Tout d'abord tu vas utiliser la fonction de récupération des resources depuis demo://content, ensuite tu vas appeler le prompt get-summary-of-requests en utilisant le contenu récupéré précedemment pour générer les attentes des participants`,
				},
			},
		},
	}, nil
}
