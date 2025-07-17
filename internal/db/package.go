package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Package struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Authors      []string          `json:"authors"`
	Dependencies map[string]string `json:"dependencies"`
	Checksum     string            `json:"checksum"`
	Filename     string            `json:"filename"`
	CreatedAt    time.Time         `json:"created_at"`
}

func InsertPackage(db *pgxpool.Pool, pkg Package) (Package, error) {
	err := db.QueryRow(
		context.Background(),
		`
            INSERT INTO packages (name, version, description, authors, dependencies, checksum, filename)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id, created_at
        `,
		pkg.Name,
		pkg.Version,
		pkg.Description,
		pkg.Authors,
		pkg.Dependencies,
		pkg.Checksum,
		pkg.Filename,
	).Scan(&pkg.ID, &pkg.CreatedAt)

	return pkg, err
}

func CheckDuplicatePackages(db *pgxpool.Pool, name, version, checksum string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM packages WHERE name = $1 AND version = $2)",
		name,
		version,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	var checksumExists bool
	err = db.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM packages WHERE checksum = $1)",
		checksum,
	).Scan(&checksumExists)
	if err != nil {
		return false, err
	}

	exists = exists && checksumExists

	return exists, err
}
