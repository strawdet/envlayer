// Package watcher provides file system watching for .env files,
// triggering a callback when any watched file changes.
package watcher

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ChangeFunc is called when a watched file is modified.
type ChangeFunc func(path string)

// Watcher monitors a set of file paths for modifications.
type Watcher struct {
	mu       sync.Mutex
	paths    map[string]time.Time
	interval time.Duration
	onChange ChangeFunc
	stop     chan struct{}
}

// New creates a new Watcher that polls at the given interval.
// onChange is invoked (in a goroutine) whenever a file's mtime changes.
func New(interval time.Duration, onChange ChangeFunc) *Watcher {
	return &Watcher{
		paths:    make(map[string]time.Time),
		interval: interval,
		onChange: onChange,
		stop:     make(chan struct{}),
	}
}

// Add registers a file path to be watched.
// If the file does not exist yet, it will be picked up once it appears.
func (w *Watcher) Add(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if info, err := os.Stat(path); err == nil {
		w.paths[path] = info.ModTime()
	} else {
		w.paths[path] = time.Time{}
	}
}

// Start begins the polling loop. Call Stop to terminate it.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.check()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop terminates the polling loop.
func (w *Watcher) Stop() {
	close(w.stop)
}

// Paths returns the currently watched file paths.
func (w *Watcher) Paths() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]string, 0, len(w.paths))
	for p := range w.paths {
		out = append(out, p)
	}
	return out
}

func (w *Watcher) check() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for path, last := range w.paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().After(last) {
			w.paths[path] = info.ModTime()
			go w.onChange(path)
		}
	}
}

// String returns a human-readable description of the watcher.
func (w *Watcher) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return fmt.Sprintf("Watcher{files: %d, interval: %s}", len(w.paths), w.interval)
}
