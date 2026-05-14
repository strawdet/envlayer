package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/envlayer/envlayer/internal/linter"
	"github.com/envlayer/envlayer/internal/loader"
	"github.com/envlayer/envlayer/internal/merger"
	"github.com/envlayer/envlayer/internal/resolver"
	"github.com/envlayer/envlayer/internal/sorter"
)

// Run is the CLI entry point.
func Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envlayer <command> [options]")
		return 1
	}

	switch args[0] {
	case "resolve":
		return runResolve(args[1:])
	case "lint":
		return runLint(args[1:])
	case "print":
		return runPrint(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		return 1
	}
}

func runResolve(args []string) int {
	fs := flag.NewFlagSet("resolve", flag.ContinueOnError)
	dir := fs.String("dir", ".", "base directory for .env files")
	ctx := fs.String("context", "", "runtime context (e.g. production)")
	format := fs.String("format", "kv", "output format: kv, export, json")
	sortOrder := fs.String("sort", "asc", "sort order: asc, desc")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	m := merger.New(*dir)
	r := resolver.New(m)
	env, err := r.Resolve(*ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve error: %v\n", err)
		return 1
	}

	ord := sorter.Ascending
	if *sortOrder == "desc" {
		ord = sorter.Descending
	}
	s := sorter.New(sorter.WithOrder(ord))
	printEnv(env, *format, s)
	return 0
}

func runLint(args []string) int {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	path := fs.String("file", ".env", "path to .env file")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	env, err := loader.LoadFile(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load error: %v\n", err)
		return 1
	}
	l := linter.New(linter.WithEmptyValueCheck(), linter.WithOSShadowCheck())
	findings := l.Lint(env)
	if len(findings) == 0 {
		fmt.Println("no issues found")
		return 0
	}
	for _, f := range findings {
		fmt.Println(f)
	}
	return 1
}

func runPrint(args []string) int {
	fs := flag.NewFlagSet("print", flag.ContinueOnError)
	path := fs.String("file", ".env", "path to .env file")
	format := fs.String("format", "kv", "output format: kv, export, json")
	sortBy := fs.String("sort-by", "key", "sort by: key, value")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	env, err := loader.LoadFile(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load error: %v\n", err)
		return 1
	}
	opts := []sorter.Option{sorter.WithCaseInsensitive()}
	if *sortBy == "value" {
		opts = append(opts, sorter.WithSortByValue())
	}
	s := sorter.New(opts...)
	printEnv(env, *format, s)
	return 0
}

func printEnv(env map[string]string, format string, s *sorter.Sorter) {
	pairs := s.Sorted(env)
	switch strings.ToLower(format) {
	case "export":
		for _, p := range pairs {
			fmt.Printf("export %s=%s\n", p[0], p[1])
		}
	case "json":
		ordered := make(map[string]string, len(pairs))
		for _, p := range pairs {
			ordered[p[0]] = p[1]
		}
		b, _ := json.MarshalIndent(ordered, "", "  ")
		fmt.Println(string(b))
	default:
		for _, p := range pairs {
			fmt.Printf("%s=%s\n", p[0], p[1])
		}
	}
}
