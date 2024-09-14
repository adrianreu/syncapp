package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var path string // Global variable for the path

func main() {
	var rootCmd = &cobra.Command{Use: "syncapp"}

	// Sync command
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "init syncapp repo",
		Run:   initCommand,
	}
	// Sync command
	var syncCmd = &cobra.Command{
		Use:   "push",
		Short: "push files for the current commit",
		Run:   syncCommand,
	}

	// Extract command
	var extractCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pulled synced files for the current branch and commit",
		Run:   extractCommand,
	}

	// Add the --path flag
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "Specify the directory to run the command")

	// Register commands
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(initCmd)

	// Execute the CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
