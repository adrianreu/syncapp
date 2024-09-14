package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/klauspost/compress/zstd"
)

// Create a Zstandard-compressed archive of the provided files and directories.
func createZstdArchive(files []string, archiveName string) error {
	// Create the output .zst file
	outFile, err := os.Create(archiveName)
	if err != nil {
		return fmt.Errorf("error creating archive file: %v", err)
	}
	defer outFile.Close()

	// Initialize the Zstandard encoder
	enc, err := zstd.NewWriter(
		outFile,
		zstd.WithEncoderConcurrency(runtime.NumCPU()),
		zstd.WithEncoderLevel(zstd.SpeedFastest),
	)
	if err != nil {
		return fmt.Errorf("error creating zstd writer: %v", err)
	}
	defer enc.Close()

	// Create a tar writer to archive multiple files/directories
	tarWriter := tar.NewWriter(enc)
	defer tarWriter.Close()

	// Iterate over the provided files and add them to the tar archive
	for _, file := range files {
		err = addFileToTar(tarWriter, file)
		if err != nil {
			return fmt.Errorf("error adding file to archive: %v", err)
		}
	}

	return nil
}

// Helper function to add a file or directory to the tar archive
func addFileToTar(tarWriter *tar.Writer, filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("could not stat file: %v", err)
	}

	// Handle directories
	if fileInfo.IsDir() {
		// Recursively add directory contents
		return filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return addTarFile(tarWriter, path, info)
		})
	}

	// Add a single file
	return addTarFile(tarWriter, filePath, fileInfo)
}

// Helper function to add an individual file or directory to the tar writer
func addTarFile(tarWriter *tar.Writer, filePath string, fileInfo os.FileInfo) error {
	// Open the file (if it's not a directory)
	var file *os.File
	var err error
	if !fileInfo.IsDir() {
		file, err = os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file %s: %v", filePath, err)
		}
		defer file.Close()
	}

	// Create a tar header based on the file's info
	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	if err != nil {
		return fmt.Errorf("error creating tar header: %v", err)
	}
	header.Name = filePath // Full path in tar archive

	// Write the tar header
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("error writing tar header: %v", err)
	}

	// Write the file contents (if it's a file)
	if !fileInfo.IsDir() {
		if _, err := io.Copy(tarWriter, file); err != nil {
			return fmt.Errorf("error writing file to tar: %v", err)
		}
	}

	return nil
}

// Extract a Zstandard-compressed archive to the specified directory.
func extractZstdArchive(archiveName, targetDir string) error {
	// Open the compressed archive file
	archiveFile, err := os.Open(archiveName)
	if err != nil {
		return fmt.Errorf("error opening archive: %v", err)
	}
	defer archiveFile.Close()

	// Initialize the Zstandard decoder
	dec, err := zstd.NewReader(archiveFile)
	if err != nil {
		return fmt.Errorf("error creating zstd reader: %v", err)
	}
	defer dec.Close()

	// Create a tar reader to extract files from the tar archive
	tarReader := tar.NewReader(dec)

	// Iterate over the tar archive entries and extract them
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			return fmt.Errorf("error reading tar entry: %v", err)
		}

		// Determine the full target path
		targetPath := filepath.Join(targetDir, header.Name)

		// Handle directories
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
			continue
		}

		// Handle regular files
		file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer file.Close()

		// Write the file contents
		if _, err := io.Copy(file, tarReader); err != nil {
			return fmt.Errorf("error writing file: %v", err)
		}
	}

	return nil
}
