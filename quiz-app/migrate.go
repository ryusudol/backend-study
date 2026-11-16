package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func migrate(db *sql.DB, dir string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version 	INTEGER 	PRIMARY KEY,
		applied_at 	TIMESTAMPTZ NOT NULL DEFAULT now()
	)`)

	if err != nil {
		return err
	}

	applied, err := appliedVersions(db)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read %s: %w", dir, err)
	}

	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		version, err := parseVersion(name)
		if err != nil {
			return err
		}

		if applied[version] {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}

		if err := applyMigration(db, version, string(content)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}

		log.Println("applied migration", name)
	}

	return nil
}

func appliedVersions(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query(`SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func parseVersion(name string) (int, error) {
	prefix, _, found := strings.Cut(name, "_")
	if !found {
		return 0, fmt.Errorf("invalid migration name %q", name)
	}

	version, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func applyMigration(db *sql.DB, version int, sqlText string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(sqlText)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, version)
	if err != nil {
		return err
	}

	return tx.Commit()
}
