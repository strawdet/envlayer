// Package archiver provides functionality to archive and restore env snapshots
// as compressed tar archives, enabling portability across environments.
package archiver

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a named env map stored in an archive.
type Entry struct {
	Name      string            `json:"name"`
	Env       map[string]string `json:"env"`
	CreatedAt time.Time         `json:"created_at"`
}

// Archiver writes and reads env archives.
type Archiver struct{}

// New returns a new Archiver.
func New() *Archiver { return &Archiver{} }

// Write serialises one or more entries into a gzip-compressed tar file at path.
func (a *Archiver) Write(path string, entries []Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("archiver: create %s: %w", path, err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, e := range entries {
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("archiver: marshal %s: %w", e.Name, err)
		}
		hdr := &tar.Header{
			Name:    e.Name + ".json",
			Size:    int64(len(data)),
			Mode:    0600,
			ModTime: e.CreatedAt,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("archiver: write header %s: %w", e.Name, err)
		}
		if _, err := tw.Write(data); err != nil {
			return fmt.Errorf("archiver: write body %s: %w", e.Name, err)
		}
	}
	return nil
}

// Read reads all entries from a gzip-compressed tar file at path.
func (a *Archiver) Read(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("archiver: open %s: %w", path, err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("archiver: gzip reader: %w", err)
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	var entries []Entry
	for {
		_, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("archiver: read header: %w", err)
		}
		var e Entry
		if err := json.NewDecoder(tr).Decode(&e); err != nil {
			return nil, fmt.Errorf("archiver: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
