package fixtures

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// NewMux registers intentionally vulnerable handlers for CodeQL Go testing.
// The fixture does not start a server or expose these handlers by itself.
func NewMux(db *sql.DB, publicDirectory string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/go/command", commandHandler)
	mux.HandleFunc("/go/sql", func(w http.ResponseWriter, r *http.Request) {
		sqlHandler(w, r, db)
	})
	mux.HandleFunc("/go/file", func(w http.ResponseWriter, r *http.Request) {
		fileHandler(w, r, publicDirectory)
	})
	mux.HandleFunc("/go/search", searchHandler)

	return mux
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")

	// INTENTIONAL VULNERABILITY: request data reaches a shell command.
	output, err := exec.Command("sh", "-c", "ping -c 1 "+host).CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(output)
}

func sqlHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	username := r.URL.Query().Get("username")

	// INTENTIONAL VULNERABILITY: untrusted input is concatenated into SQL.
	query := "SELECT id, username FROM users WHERE username = '" + username + "'"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintln(w, "query completed")
}

func fileHandler(w http.ResponseWriter, r *http.Request, publicDirectory string) {
	requestedFile := r.URL.Query().Get("file")

	// INTENTIONAL VULNERABILITY: the joined path is not confined to the root.
	contents, err := os.ReadFile(filepath.Join(publicDirectory, requestedFile))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	_, _ = w.Write(contents)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")

	// INTENTIONAL VULNERABILITY: request data is written as HTML without escaping.
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Results for %s</h1></body></html>", term)
}
