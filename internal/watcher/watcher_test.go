package watcher_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"envlayer/internal/watcher"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "FOO=bar\n")

	var mu sync.Mutex
	var changed []string

	w := watcher.New(20*time.Millisecond, func(p string) {
		mu.Lock()
		changed = append(changed, p)
		mu.Unlock()
	})
	w.Add(path)
	w.Start()
	defer w.Stop()

	// Allow watcher to record initial mtime.
	time.Sleep(30 * time.Millisecond)

	// Modify the file.
	if err := os.WriteFile(path, []byte("FOO=baz\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Wait for detection.
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(changed) == 0 {
		t.Error("expected change event, got none")
	}
	if changed[0] != path {
		t.Errorf("expected path %q, got %q", path, changed[0])
	}
}

func TestWatcher_NoFalsePositive(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "KEY=val\n")

	var mu sync.Mutex
	var count int

	w := watcher.New(20*time.Millisecond, func(_ string) {
		mu.Lock()
		count++
		mu.Unlock()
	})
	w.Add(path)
	w.Start()
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 0 {
		t.Errorf("expected 0 change events, got %d", count)
	}
}

func TestWatcher_Paths(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTempEnv(t, dir, ".env", "A=1\n")
	p2 := writeTempEnv(t, dir, ".env.local", "B=2\n")

	w := watcher.New(time.Second, func(_ string) {})
	w.Add(p1)
	w.Add(p2)

	paths := w.Paths()
	if len(paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(paths))
	}
}

func TestWatcher_String(t *testing.T) {
	w := watcher.New(500*time.Millisecond, func(_ string) {})
	s := w.String()
	if s == "" {
		t.Error("expected non-empty String()")
	}
}
