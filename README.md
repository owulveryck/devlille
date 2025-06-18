# DevLille 2025 - MCP Demo: La Révolution des API pour l'IA

Ce repository contient une démonstration interactive du Model Context Protocol (MCP) présentée au DevLille 2025. La démo illustre concrètement comment MCP permet aux IA d'interagir avec des outils et services de manière standardisée.

## Vue d'ensemble de la démonstration (40 minutes)

### Architecture de la démo

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   web_interface │    │      server      │    │  openaiserver   │
│                 │    │  (MCP Server)    │    │                 │
│  Collecte des   │────▶│                  │────▶│  Interface      │
│  attentes via   │    │  JSON-RPC via    │    │  Big-AGI        │
│  formulaire web │    │  STDIO           │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Déroulement de la démonstration

### Phase 1: Collecte des attentes (9 minutes)

1. **Lancement de l'interface web** (`demo/web_interface/`)
   - Démarre un serveur web avec NGrok
   - Génère automatiquement un QR code qui est intégré dans les slides
   - Interface simple permettant aux participants de saisir leurs attentes
   - Sauvegarde automatique dans `demo/DB/db.json`

2. **Interaction avec le public**
   - Les participants scannent le QR code
   - Saisissent leurs attentes pour la conférence
   - Les données sont collectées en temps réel

3. **Démonstration MCP en terminal**
   - Présentation du serveur MCP (`demo/server/`)
   - Commandes JSON-RPC via STDIO (exemples dans `assets/`)
   - Utilisation de `cat assets/*.json | demo/server/server` pour lister les informations
   - Explication que STDIO est le transport MCP

### Phase 2: Intégration avec l'IA (9 minutes)

4. **Lancement du serveur OpenAI**
   - Démarrage d'`openaiserver` qui utilise le même serveur MCP
   - Interface Big-AGI pour interagir avec l'IA
   - Démonstration que le serveur devient utilisable par une IA

5. **Synthèse par l'IA**
   - Demande à l'IA de synthétiser les attentes des participants
   - Illustration de la "magie" MCP : données accessibles à l'IA

### Phase 3: Modification dynamique des slides (8 minutes)

6. **Génération de contenu**
   - L'IA modifie les slides en fonction des attentes
   - Auto-reload automatique grâce au serveur principal (`main.go`)
   - Nouvelles slides générées dynamiquement dans `slides/`

### Phase 4: Explication technique (8 minutes)

7. **Architecture MCP révélée**
   - Explication des actions qui modifient les fichiers
   - Focus sur les prompts qui guident l'IA (`demo/server/prompt.go`)
   - Démonstration du code source avec les prompts intégrés

### Phase 5: Présentation des 3 systèmes MCP (6 minutes)

8. **Récapitulatif des composants**
   - Serveur MCP de collecte de données
   - Interface OpenAI/Big-AGI
   - Système de modification des slides
   - Transition vers le reste de la présentation

## Structure du projet

### Composants principaux

- **`demo/web_interface/`** - Interface web de collecte des attentes
  - Serveur HTTP avec NGrok
  - Génération automatique de QR code
  - Sauvegarde en JSON

- **`demo/server/`** - Serveur MCP principal
  - Expose les données via le protocole MCP
  - Transport JSON-RPC sur STDIO
  - Ressource `demo://content` avec les attentes
  - Prompts intégrés pour guider l'IA

- **`main.go`** - Serveur de présentation Reveal.js
  - Auto-reload des slides
  - Serveur de fichiers statiques
  - WebSocket pour le live reload

### Fichiers de configuration

- **`assets/`** - Exemples de commandes JSON-RPC MCP
- **`slides/`** - Slides individuelles de la présentation
- **`demo/DB/`** - Base de données JSON des attentes

## Utilisation

### Prérequis

- Go 1.21+
- NGrok (avec token d'authentification)
- Big-AGI ou client compatible OpenAI

### Démarrage de la démo

1. **Interface de collecte**
   ```bash
   cd demo/web_interface
   go run main.go
   ```

2. **Serveur MCP** (pour tests manuels)
   ```bash
   cd demo/server
   cat ../assets/01_init.json | go run main.go
   ```

3. **Serveur de présentation**
   ```bash
   go run main.go
   # Accessible sur https://localhost:8081
   ```

### Test du protocole MCP

Exemples de commandes disponibles dans `assets/`:
- `01_init.json` - Initialisation du protocole
- `02_tool_list.json` - Liste des outils disponibles
- `03_resources_list.json` - Liste des ressources
- `04_prompts_list.json` - Liste des prompts
- `05_resource_read.json` - Lecture du contenu des attentes

## Objectifs pédagogiques

Cette démonstration illustre concrètement :

1. **Le protocole MCP** - Communication standardisée entre IA et outils
2. **La simplicité d'intégration** - Réutilisation du même serveur MCP
3. **Les prompts comme guide** - Direction de l'IA via des instructions
4. **L'interactivité** - Adaptation du contenu en temps réel
5. **Les trois composants MCP** - Collecte, traitement, action

## Technologies utilisées

- **Go** - Langage principal
- **MCP-Go** - Implémentation Go du Model Context Protocol
- **Reveal.js** - Framework de présentation
- **NGrok** - Tunneling pour l'accès public
- **WebSocket** - Auto-reload des slides
- **Big-AGI** - Interface utilisateur pour l'IA

## Structure des données

### Format des attentes collectées

```json
[
  {
    "id": 1,
    "text": "Attente du participant",
    "user_agent": "Mozilla/5.0...",
    "timestamp": "2025-06-18T10:30:00Z"
  }
]
```

### Ressource MCP exposée

- **URI**: `demo://content`
- **Type MIME**: `text/plain`
- **Contenu**: Concaténation des attentes (limitée à 255 caractères)

## Licence

Projet de démonstration pour DevLille 2025 - Usage éducatif
