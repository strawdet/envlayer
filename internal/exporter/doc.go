// Package exporter provides functionality for writing merged environment
// variables to various output formats.
//
// Supported formats:
//
//   - dotenv  — Standard KEY="value" pairs, suitable for .env files.
//   - export  — Shell-compatible export KEY="value" statements.
//   - json    — A JSON object mapping keys to string values.
//
// Example usage:
//
//	vars := map[string]string{
//		"APP_ENV": "production",
//		"DB_HOST": "db.example.com",
//	}
//
//	e := exporter.New(vars, exporter.FormatExport)
//	if err := e.Write(os.Stdout); err != nil {
//		log.Fatal(err)
//	}
//
// Output can also be written directly to a file using WriteToFile.
package exporter
