package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func syncCommand(cmd *cobra.Command, args []string) {
	if err := os.Chdir(path); err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	if err := loadSyncConfig(); err != nil {
		fmt.Println("Error loading .syncconfig:", err)
		return
	}

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

	var matchedFiles []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matchesSyncPatternAndCheckFileSize(path) {
			matchedFiles = append(matchedFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error finding matching files:", err)
		return
	}

	if len(matchedFiles) == 0 {
		fmt.Println("No files matched the patterns.")
		return
	}

	zipName := filepath.Join(syncConfig.CloudDir, fmt.Sprintf("%s_%s.zstd", branch, commit))

	if syncConfig.KeepLatest {
		if err := removePreviousZipFiles(branch); err != nil {
			fmt.Println("Error removing old zip files:", err)
			return
		}
	}

	if err := createZstdArchive(matchedFiles, zipName); err != nil {
		fmt.Println("Error creating zip archive:", err)
		return
	}

	fmt.Println("Zip file created:", zipName)

	if err := updateGitIgnore(); err != nil {
		fmt.Println("Error updating .gitignore:", err)
		return
	}

	fmt.Println("Sync complete.")
}
