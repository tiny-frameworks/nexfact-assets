// setup.go
package setup

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"codeberg.org/tiny-frameworks/nexutils/errors"
)

const ( // provider constants
	CONTAINER = "container"
	NATIVE    = "native"
)

// SetupTestEnv copies the central testdata directory into a temporary directory in isolation.
// func SetupTestEnv(provider string) (env, resources string) {
func SetupTestEnv(provider string) (map[string]string, error) {

	startDir, _ := filepath.Abs("./")

	// Sucht nach dem Ordner "nexfact-assets"
	assetsDir, err := FindSiblingModule(startDir, "tiny-frameworks")
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

	return setupProvider(provider, testTarget)
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
func setupProvider(provider string, targetDir string) (map[string]string, error) {

	envRoot, _ := filepath.Abs(filepath.Join(targetDir, "env"))
	resourcePath, _ := filepath.Abs(filepath.Join(targetDir, "resources"))

	systemNativeResource := filepath.Join(resourcePath, "system.native")
	systemContainerResource := filepath.Join(resourcePath, "system.container")
	systemYamlTarget := filepath.Join(envRoot, "system", "system.yaml")

	switch provider {
	case NATIVE:
		copyFile(systemNativeResource, systemYamlTarget)
	case CONTAINER:
		copyFile(systemContainerResource, systemYamlTarget)
	}
	return map[string]string{
		"envRoot":      envRoot,
		"resourcePath": resourcePath,
		"tempDir":      targetDir}, nil
}

// FindSiblingModule sucht ausgehend von startDir nach oben nach einem benachbarten Ordner/Modul.
func FindSiblingModule(startDir, targetModuleName string) (string, error) {
	current := filepath.Clean(startDir)

	for {
		// Prüfe, ob targetModuleName als Ordner im aktuellen Verzeichnis existiert
		candidate := filepath.Join(current, targetModuleName)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}

		// Nächsthöheres Verzeichnis ermitteln
		parent := filepath.Dir(current)
		if parent == current {
			// Wurzelverzeichnis erreicht (/) -> nicht gefunden
			break
		}
		current = parent
	}

	return "", fmt.Errorf("module %s not found in parent directories of %s", targetModuleName, startDir)
}
