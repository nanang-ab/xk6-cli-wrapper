package amscli

import (
	"context"
	"crypto/sha256"
	"fmt"
	"gopkg.in/errgo.v2/fmt/errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/cli-wrapper", new(CLIWrapper))
}

type CLIWrapper struct {
}

func (c *CLIWrapper) DownloadCLIFile(url, dir string) (string, error) {
	// Parse the file name from the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	fileName := filepath.Base(resp.Request.URL.Path)
	filePath := filepath.Join(dir, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the content from the response to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Make the file executable
	if err := os.Chmod(filePath, 0755); err != nil {
		return "", err
	}

	return filePath, nil
}

func (c *CLIWrapper) ValidateCLIFileHash(cliPath, expectedHash string) (bool, error) {
	file, err := os.Open(cliPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false, err
	}
	computedHash := fmt.Sprintf("%x", hasher.Sum(nil))

	return computedHash == expectedHash, nil
}

// ExecuteCommand runs the specified CLIWrapper command with arguments and returns its output
func (c *CLIWrapper) ExecuteCommand(cliPath string, args ...string) (string, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, cliPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		reason := fmt.Sprintf("Command exectuion failed.\n"+
			"=> cliPath: %s\n"+
			"=> args: %#v\n"+
			"=> error status: %v\n"+
			"=> command output: \n%s", cliPath, args, err, output)
		reason = strings.TrimSpace(reason)
		return "", errors.Newf(reason)
	}
	return string(output), err
}

func (c *CLIWrapper) GetAbsolutePath(cliPath string) (string, error) {
	absPath, err := filepath.Abs(cliPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// GetWritableDirectory attempts to find a writable directory,
// defaults to the system temp directory if specified directory is not writable.
func (c *CLIWrapper) GetWritableDirectory(baseDirectory string) string {

	if isWritable(baseDirectory) {
		return baseDirectory
	}
	// Fallback to system temp directory
	return filepath.ToSlash(os.TempDir())
}

// isWritable checks if a directory is writable
func isWritable(path string) bool {
	// Attempt to create a temporary file in the directory
	testFile := filepath.Join(path, ".testwrite")
	f, err := os.Create(testFile)
	if err != nil {
		return false
	}
	// Clean up
	defer os.Remove(testFile)

	err = f.Close()
	if err != nil {
		return false
	}

	return true
}

// CleanupCLI removes the downloaded CLI tool.
func (c *CLIWrapper) CleanupCLI(cliPath string) error {
	absPath, err := filepath.Abs(cliPath)
	if err != nil {
		return err
	}

	// Perform the cleanup
	err = os.Remove(absPath)
	if err != nil {
		return err
	}

	return nil
}
