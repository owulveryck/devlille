package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Submission struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}

type Database struct {
	mu          sync.RWMutex
	submissions []Submission
	nextID      int
}

func (db *Database) Add(text, userAgent string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	submission := Submission{
		ID:        db.nextID,
		Text:      text,
		UserAgent: userAgent,
		Timestamp: time.Now(),
	}

	db.submissions = append(db.submissions, submission)
	db.nextID++
}

func (db *Database) GetAll() []Submission {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make([]Submission, len(db.submissions))
	copy(result, db.submissions)
	return result
}

func (db *Database) SaveToFile(filename string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := json.MarshalIndent(db.submissions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (db *Database) LoadFromFile(filename string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var submissions []Submission
	if err := json.Unmarshal(data, &submissions); err != nil {
		return err
	}

	db.submissions = submissions
	if len(submissions) > 0 {
		maxID := 0
		for _, s := range submissions {
			if s.ID > maxID {
				maxID = s.ID
			}
		}
		db.nextID = maxID + 1
	}

	return nil
}

var (
	db     = &Database{nextID: 1}
	dbPath = flag.String("db", "../DB/db.json", "Path to the database file")
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Text Submission</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        .form-container { background: #f5f5f5; padding: 30px; border-radius: 8px; margin-bottom: 30px; }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; resize: vertical; }
        button { background: #007cba; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; }
        button:hover { background: #005a8a; }
        .submissions { background: white; padding: 20px; border-radius: 8px; border: 1px solid #ddd; }
        .submission { margin-bottom: 15px; padding: 15px; background: #f9f9f9; border-radius: 4px; }
        .submission-meta { font-size: 12px; color: #666; margin-bottom: 8px; }
        .submission-text { font-size: 14px; }
    </style>
</head>
<body>
    <h1>Text Submission Form</h1>
    
    <div class="form-container">
        <form method="POST" action="/submit">
            <div class="form-group">
                <label for="text">Enter your text:</label>
                <textarea id="text" name="text" rows="4" placeholder="Type your message here..." required></textarea>
            </div>
            <button type="submit">Submit</button>
        </form>
    </div>

    <div class="submissions">
        <h2>Submissions ({{.Count}})</h2>
        {{range .Submissions}}
        <div class="submission">
            <div class="submission-meta">
                ID: {{.ID}} | {{.Timestamp.Format "2006-01-02 15:04:05"}} | User Agent: {{.UserAgent}}
            </div>
            <div class="submission-text">{{.Text}}</div>
        </div>
        {{else}}
        <p>No submissions yet.</p>
        {{end}}
    </div>
</body>
</html>
`

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	submissions := db.GetAll()
	data := struct {
		Submissions []Submission
		Count       int
	}{
		Submissions: submissions,
		Count:       len(submissions),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}

	db.Add(text, userAgent)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func startPeriodicSave() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			if err := db.SaveToFile(*dbPath); err != nil {
				log.Printf("Error saving database: %v", err)
			} else {
				log.Printf("Database saved to %s", *dbPath)
			}
		}
	}()
}

func main() {
	flag.Parse()

	if err := db.LoadFromFile(*dbPath); err != nil {
		log.Printf("Error loading database: %v", err)
	} else {
		log.Printf("Database loaded from %s", *dbPath)
	}

	startPeriodicSave()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)

	fmt.Println("Server starting on http://localhost:8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
