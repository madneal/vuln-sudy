package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// The handlers below intentionally keep the source-to-sink path indirect.
// Each example crosses several functions, a value object, and a service
// boundary so that interprocedural CodeQL data-flow can be exercised.

type commandInput struct {
	host string
}

type commandDecoder interface {
	Decode(*http.Request) commandInput
}

type queryCommandDecoder struct{}

func (queryCommandDecoder) Decode(r *http.Request) commandInput {
	rawHost := r.URL.Query().Get("host")
	return commandInput{host: strings.TrimSpace(rawHost)}
}

type shell interface {
	Run(string) ([]byte, error)
}

type unixShell struct{}

func (unixShell) Run(command string) ([]byte, error) {
	// INTENTIONAL VULNERABILITY: the workflow eventually reaches a shell sink.
	return exec.Command("sh", "-c", command).CombinedOutput()
}

type commandWorkflow struct {
	decoder commandDecoder
	shell   shell
}

func (w commandWorkflow) Execute(r *http.Request) ([]byte, error) {
	input := w.decoder.Decode(r)
	command := makePingCommand(input)
	return w.shell.Run(command)
}

func makePingCommand(input commandInput) string {
	// INTENTIONAL VULNERABILITY: user-controlled host is carried through a
	// value object and string transformation into the command builder.
	return strings.Join([]string{"ping", "-c", "1", input.host}, " ")
}

func complexCommandHandler(w http.ResponseWriter, r *http.Request) {
	workflow := commandWorkflow{
		decoder: queryCommandDecoder{},
		shell:   unixShell{},
	}
	output, err := workflow.Execute(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(output)
}

type searchInput struct {
	username string
}

type searchDecoder struct{}

func (searchDecoder) Decode(r *http.Request) searchInput {
	return searchInput{username: strings.TrimSpace(r.URL.Query().Get("username"))}
}

type accountRepository struct {
	db *sql.DB
}

func (repo accountRepository) Find(ctx context.Context, input searchInput) (*sql.Rows, error) {
	statement := makeUserStatement(input)
	// INTENTIONAL VULNERABILITY: a tainted statement reaches database/sql.
	return repo.db.QueryContext(ctx, statement)
}

func makeUserStatement(input searchInput) string {
	// INTENTIONAL VULNERABILITY: the value object is interpolated into SQL.
	return fmt.Sprintf(
		"SELECT id, username FROM users WHERE username = '%s'",
		input.username,
	)
}

func complexSQLHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	input := (searchDecoder{}).Decode(r)
	rows, err := (accountRepository{db: db}).Find(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintln(w, "query completed")
}

type downloadInput struct {
	segments []string
}

func decodeDownload(r *http.Request) downloadInput {
	requested := r.URL.Query().Get("file")
	return downloadInput{segments: strings.Split(requested, "/")}
}

type artifactStore struct {
	root string
}

func (store artifactStore) Open(_ context.Context, input downloadInput) (*os.File, error) {
	parts := make([]string, 0, len(input.segments)+1)
	parts = append(parts, store.root)
	parts = append(parts, input.segments...)

	// INTENTIONAL VULNERABILITY: user-controlled path segments reach os.Open.
	return os.Open(filepath.Join(parts...))
}

func complexFileHandler(w http.ResponseWriter, r *http.Request, publicDirectory string) {
	input := decodeDownload(r)
	file, err := (artifactStore{root: publicDirectory}).Open(r.Context(), input)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	_, _ = io.Copy(w, file)
}

type resultPage struct {
	prefix string
	suffix string
}

func (page resultPage) Render(term string) string {
	// INTENTIONAL VULNERABILITY: the rendered value is not HTML-escaped.
	return page.prefix + term + page.suffix
}

func decodeSearchTerm(r *http.Request) string {
	return strings.TrimSpace(r.URL.Query().Get("q"))
}

func complexSearchHandler(w http.ResponseWriter, r *http.Request) {
	page := resultPage{
		prefix: "<html><body><h1>Results for ",
		suffix: "</h1></body></html>",
	}
	output := page.Render(decodeSearchTerm(r))
	w.Header().Set("Content-Type", "text/html")
	_, _ = io.WriteString(w, output)
}
