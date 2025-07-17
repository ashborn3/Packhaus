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

	"github.com/go-chi/chi/v5"
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

	ok, err = db.CheckDuplicatePackages(cntlr.DB, metadata.Name, metadata.Version, metadata.Version)
	if err != nil {
		http.Error(w, "duplication check failed", http.StatusInternalServerError)
		return
	}
	if ok {
		http.Error(w, "package already exists", http.StatusConflict)
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

func (cntlr *controller) DownloadPackageHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	version := chi.URLParam(r, "version")

	pkg, err := db.GetPackageByNameVersion(cntlr.DB, name, version)
	if err != nil {
		http.Error(w, "package not found", http.StatusNotFound)
		return
	}

	filePath := fmt.Sprintf("./storage/%s/%s.tar.gz", name, version)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "error reading package from disk", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", pkg.Filename))
	w.Header().Set("Content-Type", "application/gzip")

	stat, _ := file.Stat()
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "error sending file", http.StatusInternalServerError)
		return
	}
}
