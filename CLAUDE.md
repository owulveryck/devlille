# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Reveal.js presentation repository for a talk about the Model Context Protocol (MCP) titled "MCP: La Révolution des API pour l'IA" (MCP: The API Revolution for AI). The presentation explains how MCP enables AI agents to interact with tools and services in a standardized way.

## Common Development Commands

### Starting the Development Server
```bash
# Use Caddy to serve the presentation with HTTPS
caddy run --config Caddyfile
```
The presentation will be available at `https://localhost:8443`

### Building and Testing
No build process is required - this is a static HTML presentation using Reveal.js CDN resources.

## Architecture and File Structure

### Core Presentation Files
- `mcp.html` - Main 40-minute presentation (complete version)
- `mcp_courte.html` - Short version of the presentation  
- `mcp_formation.html` - Training/workshop version
- `mcp_cfp.md` - Call for Papers description document

### Static Assets
- `dist/` - Reveal.js core files (CSS, JS, themes)
- `plugin/` - Reveal.js plugins (highlight, markdown, math, notes, search, zoom)
- `coorp/images/` - Corporate branding images (OCTO logos, duck mascot)
- `mcp/` - MCP-specific images (architecture diagrams)
- `photos/` - Social media QR codes for contact information

### Configuration
- `Caddyfile` - Caddy web server configuration for HTTPS serving
- Uses Caddy's internal TLS for local development

## Key Technical Details

### Reveal.js Configuration
- Uses Reveal.js framework for HTML presentations
- Includes Mermaid.js integration for sequence diagrams
- Custom OCTO theme (`dist/theme/octo.css`)
- Plugins enabled: zoom, notes, search, markdown, highlight

### Content Structure
The main presentation (`mcp.html`) covers:
- Understanding LLMs and their limitations
- Introduction to Model Context Protocol (MCP)
- The "3U" framework: Useful, Usable, Used
- MCP vs REST API comparison
- Implementation strategies and roadmap

### Development Notes
- No package.json or Node.js build process
- Self-contained HTML files with embedded styling
- Uses CDN resources for external dependencies (Mermaid.js)
- Static file serving only - no server-side processing required

## Working with This Repository

When editing presentations:
1. Modify the HTML files directly
2. Test changes by serving with Caddy
3. The presentation includes speaker notes in `<aside class="notes">` tags
4. Mermaid diagrams are embedded as text and rendered client-side

For new features or significant changes:
1. Consider the modular plugin architecture of Reveal.js
2. Follow the existing custom CSS patterns in the embedded styles
3. Maintain consistency with the OCTO corporate branding