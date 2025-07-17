package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"packhaus/internal/db"
	"packhaus/internal/middleware"
	"path/filepath"
)

type meta struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Authors      []string          `json:"authors"`
	Dependencies map[string]string `json:"dependencies"`
	Checksum     string            `json:"checksum"`
}

func (cntlr *controller) UploadPackageHandler(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(middleware.ContextKeyUserID)
	userID, ok := val.(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32mb
	if err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	metaStr := r.FormValue("metadata")
	var metadata meta
	if err := json.Unmarshal([]byte(metaStr), &metadata); err != nil {
		http.Error(w, "invalid metadata", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := fmt.Sprintf("storage/%s/%s.tar.gz", metadata.Name, metadata.Version)

	if err := os.MkdirAll(
		filepath.Dir(path),
		0755,
	); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		fmt.Printf("error; %s\n", err.Error())
		return
	}

	out, err := os.Create(path)
	if err != nil {
		http.Error(w, "could not save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	pkg := db.Package{
		Name:         metadata.Name,
		Version:      metadata.Version,
		Description:  metadata.Description,
		Authors:      metadata.Authors,
		Dependencies: metadata.Dependencies,
		Checksum:     metadata.Checksum,
		Filename:     path,
	}

	_, err = db.InsertPackage(cntlr.DB, pkg)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
