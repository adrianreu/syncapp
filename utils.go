package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Convert a wildcard pattern to a regular expression
func wildcardToRegex(pattern string) string {
	pattern = strings.ReplaceAll(pattern, ".", "\\.") // Escape dots
	pattern = strings.ReplaceAll(pattern, "*", ".*")   // Replace * with .*
	return "^" + pattern + "$"                         // Match the whole string
}

func matchesSyncPattern(filePath string) bool {
	for _, pattern := range syncConfig.Patterns {
		// Convert wildcard pattern to regex
		regexPattern := wildcardToRegex(pattern)
		matched, _ := regexp.MatchString(regexPattern, filePath)
		if matched {
			return true
		}
	}
	return false
}
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
