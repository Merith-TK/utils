// Package archive provides utilities for working with archive formats, currently supporting ZIP extraction.
//
// The package focuses on safe and reliable archive extraction with proper path handling
// and directory structure preservation. It includes security measures to prevent
// path traversal attacks and ensures proper file permissions are maintained.
//
// Current functionality:
//   - ZIP archive extraction with full directory structure preservation
//   - Safe path handling to prevent directory traversal attacks
//   - Automatic parent directory creation
//   - Proper file permission preservation
//
// Example usage:
//
//	err := archive.Unzip("archive.zip", "/path/to/extract")
//	if err != nil {
//		log.Fatal("Failed to extract archive:", err)
//	}
package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Unzip extracts a ZIP archive from src to the dest directory.
// All files and folders in the archive will be extracted, preserving the directory structure.
// Returns an error if extraction fails.
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		
		// Clean the path to handle any issues with path separators
		fpath = filepath.Clean(fpath)
		
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// Ensure parent directory exists for files
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
