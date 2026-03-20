package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract unzips the given zipPath to destDir securely.
func Extract(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)) {
			// prevent zip slip vulnerability
			continue
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(target), 0755)
		out, err := os.Create(target)
		if err != nil {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			out.Close()
			continue
		}
		io.Copy(out, rc)
		rc.Close()
		out.Close()
	}
	return nil
}
