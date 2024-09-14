package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Create sync.yaml and setup Git hooks if not already present
func initCommand(cmd *cobra.Command, args []string) {
	// 1. Create sync.yaml if it doesn't exist
	if _, err := os.Stat("sync.yaml"); os.IsNotExist(err) {
		defaultConfig := `# sync.yaml - SyncApp configuration file
cloud_dir: "/path/to/cloud"
keep_latest: true
patterns:
  - "assets/*"
`
		err := ioutil.WriteFile("sync.yaml", []byte(defaultConfig), 0644)
		if err != nil {
			fmt.Errorf("error creating sync.yaml: %v", err)
		}
		fmt.Println("sync.yaml created with default configuration.")
	} else {
		fmt.Println("sync.yaml already exists.")
	}

	// 2. Set up Git hooks for pre-push and post-checkout
	if err := setupGitHooks(); err != nil {
		fmt.Errorf("error setting up git hooks: %v", err)
	}

	fmt.Println("Git hooks setup successfully.")
	fmt.Println("Please edit sync.yaml with your configuration.")
}

// Setup pre-push and post-checkout Git hooks to run syncapp
func setupGitHooks() error {
	// Check if .git/hooks exists
	hookDir := filepath.Join(".git", "hooks")
	if _, err := os.Stat(hookDir); os.IsNotExist(err) {
		return fmt.Errorf(".git/hooks directory does not exist, make sure this is a Git repository")
	}

	// 1. Pre-push hook to run syncapp push
	prePushHookPath := filepath.Join(hookDir, "pre-push")
	prePushHook := "#!/bin/sh\n# Pre-push hook for SyncApp\nsyncapp push\n"
	if err := writeHook(prePushHookPath, prePushHook); err != nil {
		return err
	}

	// 2. Post-checkout hook to run syncapp pull
	postCheckoutHookPath := filepath.Join(hookDir, "post-checkout")
	postCheckoutHook := "#!/bin/sh\n# Post-checkout hook for SyncApp\nsyncapp pull\n"
	if err := writeHook(postCheckoutHookPath, postCheckoutHook); err != nil {
		return err
	}

	return nil
}

// Write the Git hook file, making sure it is executable
func writeHook(hookPath, hookContent string) error {
	if err := ioutil.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("Error writing Git hook %s: %v", hookPath, err)
	}
	fmt.Printf("Git hook created: %s\n", hookPath)
	return nil
}
