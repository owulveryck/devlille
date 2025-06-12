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
		mcp.WithPromptDescription("Ce prompt permet d'adapter les slides."),
	)
}

func GetPromptHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Ce prompt permet de mettre à jour les slides",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: `
Tu es un assistant spécialisé dans la création de présentations. Tu vas lire plusieurs fichiers de présentation et les synthétiser pour répondre aux attentes des participants.

1.  Lis le contenu du fichier /Users/olivierwulveryck/github.com/owulveryck/devlille/slides/slides.txt. Chaque ligne de ce fichier représente le nom d'un fichier de présentation situé dans le répertoire /Users/olivierwulveryck/github.com/owulveryck/devlille/slides/.

2.  Lis le contenu de chaque fichier de présentation listé dans /Users/olivierwulveryck/github.com/owulveryck/devlille/slides/slides.txt. Attention le numéro du fichier ne correspond pas au numéro du slide, c'est le numéro de ligne dans le fichier qui donne le numéro du slide.

3.  Analyse et synthétise le contenu de toutes les présentations. Identifie les parties de la présentation qui répondent aux attentes potentielles des participants.

4.  Génère un contenu HTML pour une slide reveal.js qui présente les attentes des participants et indique clairement quelle partie de la présentation répond à chaque attente. Utilise des balises <section> pour encadrer le contenu de la slide.  Pour chaque attente, crée une sous-section. Par exemple:

    <section>
        <section> Attente 1 ... </section>
        <section> Attente 2 ... </section>
    </section>

    *   Utilise des emphases (gras, italique) et des emojis pour rendre la slide plus attrayante et lisible.
    *   Mets en évidence la correspondance entre les parties de la présentation et les attentes des participants.
    *   Structure la slide de manière claire et concise pour faciliter la compréhension.

5.  Remplace le contenu du fichier /Users/olivierwulveryck/github.com/owulveryck/devlille/slides/05-expectations.html par le contenu HTML généré à l'étape 4. Assure-toi que le nouveau contenu est bien formaté et valide pour reveal.js.
`,
				},
			},
		},
	}, nil
}
