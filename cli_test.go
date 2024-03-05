package amscli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const downloadUrl = "https://cdn.development.armada.accelbyte.io/linux_amd64/ams"
const cliDirectory = ".run/files/ams/cli/"

func TestDownloadAndExecute(t *testing.T) {
	cli := CLIWrapper{}

	writableCliDir := cli.GetWritableDirectory(cliDirectory)
	cliPath, err := cli.DownloadCLIFile(downloadUrl, writableCliDir)
	if err != nil {
		t.Fatalf("Failed to download: %v", err)
	}

	// Ensure the file gets cleaned up after the test
	defer os.Remove(cliPath)

	//TODO: uncomment this after we have an official hash sum
	//isValid, err := cli.ValidateCLIFileHash(cliPath, expectedHash)
	//if err != nil {
	//	t.Fatalf("Failed to validate AMS CLIWrapper %v", err)
	//}
	//if !isValid {
	//	t.Fatalf("Invalid AMS CLIWrapper hash")
	//}

	// Execute a command with the downloaded AMS CLIWrapper app
	output, err := cli.ExecuteCommand(cliPath, "--help")
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	expectedOutput := "version: "
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected command output should start with %q, got %q", expectedOutput, output)
	}
}

func TestDownloadAndGetAbsolutePath(t *testing.T) {
	cli := CLIWrapper{}

	cliPath, err := cli.DownloadCLIFile(downloadUrl, cliDirectory)
	if err != nil {
		t.Fatalf("Failed to download: %v", err)
	}

	// Ensure the file gets cleaned up after the test
	defer os.Remove(cliPath)

	// Get absolute path
	absPath, err := cli.GetAbsolutePath(cliPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	if !filepath.IsAbs(absPath) {
		t.Errorf("Expected cliPath is not absolute path")
	}
}
