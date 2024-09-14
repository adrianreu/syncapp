package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// TODO fix this
func extractCommand(cmd *cobra.Command, args []string) {
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

	zipName := filepath.Join(syncConfig.CloudDir, fmt.Sprintf("%s_%s.zstd", branch, commit))

	if err := extractZstdArchive(zipName, "."); err != nil {
		fmt.Println("Error extracting zip archive:", err)
		return
	}

	fmt.Println("Files extracted to:", path)
}
