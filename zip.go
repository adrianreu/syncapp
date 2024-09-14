package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Create a zip archive for files and directories
func createZipArchive(files []string, archiveName string) error {
	zipFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		err = addToZip(zipWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// Recursively add files and directories to the zip archive
func addToZip(zipWriter *zip.Writer, filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return addDirectoryToZip(zipWriter, filePath, "")
	}

	// Add a file to the zip
	return addFileToZip(zipWriter, filePath)
}

// Add a file to the zip archive
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filePath
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// Add a directory and its contents to the zip archive
func addDirectoryToZip(zipWriter *zip.Writer, dirPath string, prefix string) error {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			// Recursively add directories
			if err := addDirectoryToZip(zipWriter, fullPath, filepath.Join(prefix, entry.Name())); err != nil {
				return err
			}
		} else {
			// Add files
			if err := addFileToZip(zipWriter, fullPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractZip(zipName string, extractDir string) error {
	zipFile, err := zip.OpenReader(zipName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		fpath := filepath.Join(extractDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		if _, err := io.Copy(dstFile, rc); err != nil {
			return err
		}
	}
	return nil
}
