package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type SyncConfig struct {
	CloudDir   string   `yaml:"cloud_dir"`
	KeepLatest bool     `yaml:"keep_latest"`
	Patterns   []string `yaml:"patterns"`
}

var syncConfig SyncConfig

// Load `.syncconfig` file
func loadSyncConfig() error {
	data, err := ioutil.ReadFile(".syncconfig")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &syncConfig)
}

// Get the current git branch
func getGitBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Get the current git commit short hash
func getGitCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Check if a file matches any pattern in the config
func matchesSyncPattern(filePath string) bool {
	for _, pattern := range syncConfig.Patterns {
		matched, _ := regexp.MatchString(pattern, filePath)
		if matched {
			return true
		}
	}
	return false
}

// Create a zip archive for files matching the patterns
func createZipArchive(files []string, archiveName string) error {
	zipFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		err = addFileToZip(zipWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
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

// Update .gitignore to exclude synced files
func updateGitIgnore(files []string) error {
	gitIgnore, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return err
	}

	ignoreContent := string(gitIgnore)
	for _, file := range files {
		if !strings.Contains(ignoreContent, file) {
			ignoreContent += "\n" + file
		}
	}

	return ioutil.WriteFile(".gitignore", []byte(ignoreContent), 0644)
}

// Remove previous zip files of the same branch
func removePreviousZipFiles(branch string) error {
	files, err := ioutil.ReadDir(syncConfig.CloudDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), branch+"_") && strings.HasSuffix(file.Name(), ".zip") {
			err := os.Remove(filepath.Join(syncConfig.CloudDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Check for changes in files by comparing their hashes
func hasChanges(files []string) (bool, error) {
	var oldHash, newHash []byte
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return false, err
		}
		hash := sha256.Sum256(content)
		newHash = append(newHash, hash[:]...)
	}

	hashFile := ".lasthash"
	if _, err := os.Stat(hashFile); err == nil {
		oldHash, err = ioutil.ReadFile(hashFile)
		if err != nil {
			return false, err
		}
	}

	if bytes.Equal(oldHash, newHash) {
		return false, nil
	}

	return ioutil.WriteFile(hashFile, newHash, 0644) == nil, nil
}

func main() {
	// Load configuration from .syncconfig
	if err := loadSyncConfig(); err != nil {
		fmt.Println("Error loading .syncconfig:", err)
		return
	}

	// Get Git branch and commit info
	branch, err := getGitBranch()
	if err != nil {
		fmt.Println("Error getting git branch:", err)
		return
	}

	commit, err := getGitCommit()
	if err != nil {
		fmt.Println("Error getting git commit:", err)
		return
	}

	// Find files matching patterns
	var matchedFiles []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && matchesSyncPattern(path) {
			matchedFiles = append(matchedFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error finding matching files:", err)
		return
	}

	// Check for file changes
	changed, err := hasChanges(matchedFiles)
	if err != nil {
		fmt.Println("Error checking file changes:", err)
		return
	}

	// Only create zip if changes detected
	if changed {
		zipName := filepath.Join(syncConfig.CloudDir, fmt.Sprintf("%s_%s.zip", branch, commit))

		// Remove old zip if keepLatest is true
		if syncConfig.KeepLatest {
			if err := removePreviousZipFiles(branch); err != nil {
				fmt.Println("Error removing old zip files:", err)
				return
			}
		}

		// Create the zip archive
		if err := createZipArchive(matchedFiles, zipName); err != nil {
			fmt.Println("Error creating zip archive:", err)
			return
		}

		fmt.Println("Zip file created:", zipName)
	} else {
		fmt.Println("No changes detected. Renaming zip file.")
		// Rename zip if no changes but commit changed
		oldZipName := filepath.Join(syncConfig.CloudDir, fmt.Sprintf("%s_*.zip", branch))
		newZipName := filepath.Join(syncConfig.CloudDir, fmt.Sprintf("%s_%s.zip", branch, commit))
		os.Rename(oldZipName, newZipName)
	}

	// Update .gitignore
	if err := updateGitIgnore(matchedFiles); err != nil {
		fmt.Println("Error updating .gitignore:", err)
		return
	}

	fmt.Println("Sync complete.")
}
