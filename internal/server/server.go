package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/frostyard/site/internal/build"
	"github.com/fsnotify/fsnotify"
)

// Config holds the dev server configuration.
type Config struct {
	ContentDir string
	StaticDir  string
	OutputDir  string
	Addr       string
	Root       string
}

const liveReloadScript = `<script>
const es = new EventSource('/_reload');
es.onmessage = () => location.reload();
es.onerror = () => setTimeout(() => location.reload(), 1000);
</script>`

// Serve starts the development server with file watching and live reload.
func Serve(cfg Config) error {
	// Run initial build
	buildCfg := build.Config{
		ContentDir: cfg.ContentDir,
		StaticDir:  cfg.StaticDir,
		OutputDir:  cfg.OutputDir,
		Root:       cfg.Root,
	}
	if err := build.Build(buildCfg); err != nil {
		return fmt.Errorf("initial build failed: %w", err)
	}

	// SSE client tracking
	var (
		mu      sync.Mutex
		clients []chan struct{}
	)

	addClient := func() chan struct{} {
		mu.Lock()
		defer mu.Unlock()
		ch := make(chan struct{}, 1)
		clients = append(clients, ch)
		return ch
	}

	removeClient := func(ch chan struct{}) {
		mu.Lock()
		defer mu.Unlock()
		for i, c := range clients {
			if c == ch {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}

	notifyClients := func() {
		mu.Lock()
		defer mu.Unlock()
		for _, ch := range clients {
			select {
			case ch <- struct{}{}:
			default:
			}
		}
	}

	// Set up fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}
	defer watcher.Close()

	// Watch directories recursively
	watchDirs := []string{cfg.ContentDir, cfg.StaticDir}
	templatesDir := filepath.Join(cfg.Root, "templates")
	if _, err := os.Stat(templatesDir); err == nil {
		watchDirs = append(watchDirs, templatesDir)
	}

	for _, dir := range watchDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("watching %s: %w", dir, err)
		}
	}

	// Start file watcher goroutine with debounce
	go func() {
		var debounceTimer *time.Timer

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) == 0 {
					continue
				}

				// If a new directory was created, watch it too
				if event.Op&fsnotify.Create != 0 {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						_ = watcher.Add(event.Name)
					}
				}

				// Debounce: reset timer on each event
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(200*time.Millisecond, func() {
					fmt.Println("Change detected, rebuilding...")
					if err := build.Build(buildCfg); err != nil {
						fmt.Fprintf(os.Stderr, "Rebuild failed: %v\n", err)
						return
					}
					notifyClients()
				})

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
			}
		}
	}()

	// Set up HTTP handlers with a custom ServeMux
	mux := http.NewServeMux()

	// SSE endpoint for live reload
	mux.HandleFunc("/_reload", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ch := addClient()
		defer removeClient(ch)

		// Send initial comment to establish connection
		fmt.Fprintf(w, ": connected\n\n")
		flusher.Flush()

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				fmt.Fprintf(w, "data: reload\n\n")
				flusher.Flush()
			}
		}
	})

	// File server with live reload script injection
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path

		// Determine file path
		filePath := filepath.Join(cfg.OutputDir, filepath.Clean(urlPath))

		// If path is a directory, serve index.html
		info, err := os.Stat(filePath)
		if err == nil && info.IsDir() {
			filePath = filepath.Join(filePath, "index.html")
		}

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// For HTML files, inject live reload script
		if strings.HasSuffix(filePath, ".html") {
			data, err := os.ReadFile(filePath)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}

			// Inject live reload script before </body>
			content := bytes.Replace(data, []byte("</body>"), []byte(liveReloadScript+"\n</body>"), 1)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			io.Copy(w, bytes.NewReader(content))
			return
		}

		// For non-HTML files, serve directly
		http.ServeFile(w, r, filePath)
	})

	fmt.Printf("Dev server running at http://localhost%s\n", cfg.Addr)
	fmt.Println("Watching for changes...")

	return http.ListenAndServe(cfg.Addr, mux)
}
