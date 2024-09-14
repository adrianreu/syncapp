package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

// Create a zstd archive for files and directories
func createZstdArchive(files []string, archiveName string) error {
	zstdFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer zstdFile.Close()

	zstdWriter, _ := zstd.NewWriter(zstdFile)
	defer zstdWriter.Close()

	for _, file := range files {
		if err := addToZstd(zstdWriter, file, ""); err != nil {
			return err
		}
	}
	return nil
}

// Recursively add files and directories to the zstd archive
func addToZstd(zstdWriter *zstd.Encoder, filePath string, prefix string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return addDirectoryToZstd(zstdWriter, filePath, prefix)
	}

	// Add a file to the zstd archive
	return addFileToZstd(zstdWriter, filePath, prefix)
}

// Add a file to the zstd archive
func addFileToZstd(zstdWriter *zstd.Encoder, filePath string, prefix string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	relPath := filepath.Join(prefix, filePath)
	_, err = zstdWriter.Write([]byte(relPath + "\n")) // Store file path in the stream
	if err != nil {
		return err
	}

	_, err = io.Copy(zstdWriter, file)
	return err
}

// Add a directory and its contents to the zstd archive
func addDirectoryToZstd(zstdWriter *zstd.Encoder, dirPath string, prefix string) error {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			if err := addDirectoryToZstd(zstdWriter, fullPath, filepath.Join(prefix, entry.Name())); err != nil {
				return err
			}
		} else {
			if err := addFileToZstd(zstdWriter, fullPath, prefix); err != nil {
				return err
			}
		}
	}
	return nil
}

// Extract a zstd archive to the specified directory
func extractZstdArchive(archiveName, targetDir string) error {
	zstdFile, err := os.Open(archiveName)
	if err != nil {
		return err
	}
	defer zstdFile.Close()

	zstdReader, err := zstd.NewReader(zstdFile)
	if err != nil {
		return err
	}
	defer zstdReader.Close()

	return extractFiles(zstdReader, targetDir)
}

// Extract files from the zstd reader to the target directory
func extractFiles(reader io.Reader, targetDir string) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		filePath := scanner.Text()
		if err := createFileFromStream(reader, targetDir, filePath); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// Create files and directories from the zstd stream
func createFileFromStream(reader io.Reader, targetDir, filePath string) error {
	fullPath := filepath.Join(targetDir, filePath)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}
