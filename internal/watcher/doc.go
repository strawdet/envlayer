// Package watcher implements lightweight polling-based file watching
// for envlayer's .env file hierarchy.
//
// # Overview
//
// Rather than relying on OS-specific inotify/kqueue APIs, watcher uses
// a simple ticker to periodically stat each registered path and compare
// modification times. This keeps the implementation portable and free of
// external dependencies.
//
// # Usage
//
//	w := watcher.New(500*time.Millisecond, func(path string) {
//		fmt.Println("changed:", path)
//	})
//	w.Add(".env")
//	w.Add(".env.local")
//	w.Start()
//	defer w.Stop()
//
// The ChangeFunc is invoked in its own goroutine so callers should
// synchronise any shared state accessed inside it.
package watcher
