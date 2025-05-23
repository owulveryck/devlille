package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

// SlideData represents a single slide
type SlideData struct {
	Content template.HTML
}

// PresentationData contains all slides for the presentation
type PresentationData struct {
	Slides []SlideData
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// WebSocket connections management
var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

// loadTemplate reads the HTML template from file
func loadTemplate() (string, error) {
	content, err := ioutil.ReadFile("mcp.html.template")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// handleWebSocket handles WebSocket connections for auto-reload
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientCount := len(clients)
	clientsMu.Unlock()
	
	log.Printf("🔌 New WebSocket client connected (total: %d clients)", clientCount)

	// Clean up on disconnect
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		remainingClients := len(clients)
		clientsMu.Unlock()
		log.Printf("📱 WebSocket client disconnected (remaining: %d clients)", remainingClients)
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// broadcastReload sends reload signal to all connected WebSocket clients
func broadcastReload() {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	log.Printf("🔄 Triggering page reload for %d connected clients", len(clients))
	
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			log.Printf("Error sending reload message to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// startFileWatcher monitors files for changes and triggers reload
func startFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error creating file watcher: %v", err)
		return
	}

	log.Printf("📁 Setting up file watcher...")

	// Start the goroutine first
	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				
				// Log ALL events for debugging
				log.Printf("🔍 Raw file event: %s (op: %s)", event.Name, event.Op.String())
				
				// Watch for Write, Create, Remove, and Rename operations
				if event.Op&fsnotify.Write == fsnotify.Write || 
				   event.Op&fsnotify.Create == fsnotify.Create ||
				   event.Op&fsnotify.Remove == fsnotify.Remove ||
				   event.Op&fsnotify.Rename == fsnotify.Rename {
					
					// Skip temporary files and system files, but be more permissive
					fileName := filepath.Base(event.Name)
					if strings.Contains(fileName, ".tmp") || 
					   strings.Contains(fileName, "~") ||
					   (strings.HasPrefix(fileName, ".") && !strings.HasSuffix(fileName, ".html")) {
						log.Printf("⏭️  Skipping temporary/system file: %s", event.Name)
						continue
					}
					
					log.Printf("📝 File changed: %s (operation: %s)", event.Name, event.Op.String())
					broadcastReload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("❌ File watcher error: %v", err)
			}
		}
	}()

	// Add files to watch
	filesToWatch := []string{
		"mcp.html.template",
		"slides/slides.txt",
	}
	
	for _, file := range filesToWatch {
		err = watcher.Add(file)
		if err != nil {
			log.Printf("Error watching file %s: %v", file, err)
		} else {
			log.Printf("👀 Watching file: %s", file)
		}
	}

	// Watch slides directory for changes
	err = watcher.Add("slides")
	if err != nil {
		log.Printf("Error watching slides directory: %v", err)
	} else {
		log.Printf("👀 Watching directory: slides/")
	}

	// Watch for HTML files in slides directory
	filepath.Walk("slides", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			err = watcher.Add(path)
			if err != nil {
				log.Printf("Error watching slide file %s: %v", path, err)
			} else {
				log.Printf("👀 Watching slide: %s", path)
			}
		}
		return nil
	})

	log.Printf("✅ File watcher initialized - auto-reload ready!")
}

// readSlidesFromFile reads the list of slide files from slides.txt
func readSlidesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var slides []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasSuffix(line, ".json") {
			slides = append(slides, line)
		}
	}

	return slides, scanner.Err()
}

// loadSlideContent reads the content of a single slide file
func loadSlideContent(slideFile string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join("slides", slideFile))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// loadAllSlides loads all slides based on the slides.txt file
func loadAllSlides() (PresentationData, error) {
	slideFiles, err := readSlidesFromFile("slides/slides.txt")
	if err != nil {
		return PresentationData{}, fmt.Errorf("error reading slides list: %v", err)
	}

	var slides []SlideData
	for _, slideFile := range slideFiles {
		content, err := loadSlideContent(slideFile)
		if err != nil {
			log.Printf("Warning: could not load slide %s: %v", slideFile, err)
			continue
		}
		slides = append(slides, SlideData{Content: template.HTML(content)})
	}

	return PresentationData{Slides: slides}, nil
}

// servePresentation handles HTTP requests and serves the generated presentation
func servePresentation(w http.ResponseWriter, r *http.Request) {
	log.Printf("📄 Serving presentation to %s", r.RemoteAddr)
	
	// Load slides data
	data, err := loadAllSlides()
	if err != nil {
		log.Printf("❌ Error loading slides: %v", err)
		http.Error(w, fmt.Sprintf("Error loading slides: %v", err), http.StatusInternalServerError)
		return
	}
	
	log.Printf("✅ Loaded %d slides successfully", len(data.Slides))

	// Load template from file
	htmlTemplate, err := loadTemplate()
	if err != nil {
		log.Printf("❌ Error loading template: %v", err)
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse and execute template
	tmpl, err := template.New("presentation").Parse(htmlTemplate)
	if err != nil {
		log.Printf("❌ Error parsing template: %v", err)
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("❌ Error executing template: %v", err)
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Start file watcher
	startFileWatcher()

	// Serve static files (CSS, JS, images)
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))
	http.Handle("/plugin/", http.StripPrefix("/plugin/", http.FileServer(http.Dir("plugin"))))
	http.Handle("/coorp/", http.StripPrefix("/coorp/", http.FileServer(http.Dir("coorp"))))
	http.Handle("/mcp/", http.StripPrefix("/mcp/", http.FileServer(http.Dir("mcp"))))
	http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("photos"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// WebSocket endpoint for auto-reload
	http.HandleFunc("/ws", handleWebSocket)

	// Serve files from current directory with custom handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If requesting the root path, serve the presentation
		if r.URL.Path == "/" {
			servePresentation(w, r)
			return
		}
		
		// For other paths, try to serve as static files from current directory
		filePath := strings.TrimPrefix(r.URL.Path, "/")
		if _, err := os.Stat(filePath); err == nil {
			http.ServeFile(w, r, filePath)
			return
		}
		
		// If file doesn't exist, serve the presentation (for clean URLs)
		servePresentation(w, r)
	})

	port := ":8081"
	log.Printf("🚀 Starting presentation server on port %s", port)
	log.Printf("🌐 Open http://localhost%s in your browser", port)
	log.Printf("🔄 Auto-reload enabled - server will watch for file changes")
	log.Printf("💡 Make changes to slides in the 'slides/' directory to see live updates")
	log.Fatal(http.ListenAndServe(port, nil))
}