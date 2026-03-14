package migrations

import (
	"embed"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
)

//go:embed *.sql
var files embed.FS

func Apply(db *sqlx.DB) error {
	entries, err := files.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migration directory: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		content, err := files.ReadFile(filepath.Clean(name))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
	}

	return nil
}
