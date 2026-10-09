// Copyright 2026 Georg Hagn (tiny-frameworks)
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/tiny-frameworks/nexutils/errors"
)

// containerRunner kapselt das Ausführen eines ephemeral Containers.
type containerRunner struct {
	imageName     string
	containerName string
	mode          string
	watchDir      string
	inDir         string
	artefactsDir  string
	logsDir       string
	timeOut       time.Duration
	runnerDir     string
}

const ( // provider constants
	container = "container"
	native    = "native"
)

// SetupTestEnv copies the central testdata directory into a temporary directory in isolation.
// func SetupTestEnv(provider string) (env, resources string) {
func SetupTestEnv(provider string) (map[string]string, error) {

	startDir, _ := filepath.Abs("./")

	// searches for folder "nexfact-assets"
	assetsDir, err := findSiblingModule(startDir, "tiny-frameworks")
	if err != nil {
		panic(err)
	}

	testSource := filepath.Join(assetsDir, "nexfact-assets", "testdata")

	// 2. Isolated target folder in the Go TempDir
	tempDir, err := os.MkdirTemp("", "nexfact-assets")
	if err != nil {
		err := errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("can not create tempdir"),
			"nexfact-assets.SetupTestEnv",
			err,
		)
		return nil, err
	}
	testTarget := filepath.Join(tempDir, "testdata")

	// 3. Copy directory
	if err := copyDir(testSource, testTarget); err != nil {
		err := errors.Wrap(
			errors.WriteError,
			fmt.Sprintf("failed to copy test env: %v", err),
			"nexfact-assets.SetupTestEnv",
			err,
		)
		return nil, err
	}

	return setupProvider(provider, testTarget, tempDir)
}

// Helper function: Finds a free port for isolated tests.
func GetFreeAddr() (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		err := errors.Wrap(
			errors.InternalError,
			fmt.Sprintf("failed to listen on free port: %v", err),
			"nexfact-assets.SetupTestEnv",
			err,
		)
		return "", err
	}
	defer l.Close()
	return l.Addr().String(), nil
}

// Helper function: Recursive directory copying
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// Helper function: Copy system.yaml based on 'provider'
func setupProvider(provider string, targetDir string, tempDir string) (map[string]string, error) {

	envRoot, _ := filepath.Abs(filepath.Join(targetDir, "env"))
	resourcePath, _ := filepath.Abs(filepath.Join(targetDir, "resources"))

	systemNativeResource := filepath.Join(resourcePath, "system.native")
	systemContainerResource := filepath.Join(resourcePath, "system.container")
	systemYamlTarget := filepath.Join(envRoot, "system", "system.yaml")

	switch provider {
	case native:
		copyFile(systemNativeResource, systemYamlTarget)
	case container:
		copyFile(systemContainerResource, systemYamlTarget)
		if err := setupContainer(envRoot); err != nil {
			return nil, err
		}
	}

	return map[string]string{
		"envRoot":      envRoot,
		"resourcePath": resourcePath,
		"tempDir":      tempDir,
	}, nil
}

// FindSiblingModule searches upwards from startDir for a neighboring folder/module.
func findSiblingModule(startDir, targetModuleName string) (string, error) {
	current := filepath.Clean(startDir)

	for {
		// Check whether targetModuleName exists as a folder in the current directory.
		candidate := filepath.Join(current, targetModuleName)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}

		// Determine the next higher directory
		parent := filepath.Dir(current)
		if parent == current {
			// Root directory reached (/) -> not found
			break
		}
		current = parent
	}

	return "", fmt.Errorf("module %s not found in parent directories of %s", targetModuleName, startDir)
}

func setupContainer(envRoot string) error {

	// Deletes the two containers nexfact-zfwatcher and nexfact-pdfwatcher (if they exist)
	ForceRemoveWatcherContainer()

	cmdDir := filepath.Join(envRoot, "sellers", "testSeller")

	pdfRunner := &containerRunner{
		imageName:     "ghcr.io/tiny-frameworks/nexfact-pdfmanager:latest",
		containerName: "nexfact-pdfwatcher",
		mode:          "watch",
		watchDir:      "/data/in",
		inDir:         filepath.Join(cmdDir, "/data/pdfs/in"),
		artefactsDir:  filepath.Join(cmdDir, "/data/artefacts"),
		logsDir:       filepath.Join(cmdDir, "/data/pdfs/logs"),
		timeOut:       130 * time.Second,
	}

	if err := startWatcher(pdfRunner, cmdDir); err != nil {
		return err
	}

	ferdRunner := &containerRunner{
		imageName:     "ghcr.io/tiny-frameworks/nexfact-zugferd:latest",
		containerName: "nexfact-zfwatcher",
		mode:          "watch",
		watchDir:      "/data/in",
		inDir:         filepath.Join(cmdDir, "/data/zugferds/in"),
		artefactsDir:  filepath.Join(cmdDir, "/data/artefacts"),
		logsDir:       filepath.Join(cmdDir, "/data/zugferds/logs"),
		timeOut:       130 * time.Second,
	}

	if err := startWatcher(ferdRunner, cmdDir); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return nil
}

func startWatcher(runner *containerRunner, cmdDir string) error {

	// 1. Start container in the background(-d) mit --rm
	startArgs := []string{
		"run",
		"-d", // Detached Mode (returns immediately)
		"--name", runner.containerName,
		"-e", fmt.Sprintf("NEXFACT_MODE=%s", runner.mode),
		"-e", fmt.Sprintf("NEXFACT_WATCH=%s", runner.watchDir),
		"-v", fmt.Sprintf("%s:/data/in:Z", runner.inDir),
		"-v", fmt.Sprintf("%s:/data/artefacts:Z,shared", runner.artefactsDir),
		"-v", fmt.Sprintf("%s:/data/logs:Z", runner.logsDir),
		runner.imageName,
	}

	cmd := exec.Command("podman", startArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to start container: %v\nOutput: %s", err, string(out))
		return fmt.Errorf("failed to start container: %v\nOutput: %s", err, string(out))
	}

	return nil

}

func IsContainerRunning(containerName string) bool {
	// returns immediately; formatting directly returns "true" or "false"
	cmd := exec.Command("podman", "container", "inspect",
		"--format", "{{.State.Running}}",
		containerName,
	)

	output, err := cmd.Output()
	if err != nil {
		// If the command fails, the container usually does not exist at all.
		fmt.Printf("container %s failed ", containerName)
		return false
	}

	return strings.TrimSpace(string(output)) == "true"
}

func ForceRemoveWatcherContainer() {
	// Guarantees the deletion of the containers, regardless of their state.
	containers := []string{"nexfact-zfwatcher", "nexfact-pdfwatcher"}

	for _, container := range containers {
		_ = exec.Command("podman", "rm", "-f", container).Run()
	}

}
