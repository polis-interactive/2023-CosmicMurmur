package api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func checkIfIsHtmlRoot(filePath string, rootDir string) (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	if !strings.Contains(basePath, rootDir) {
		return "", errors.New(fmt.Sprintf("PATH DOES NOT INCLUDE PROGRAM %s", rootDir))
	}
	rootPath := strings.Split(basePath, rootDir)[0]
	htmlPath := filepath.Join(rootPath, rootDir, filePath)
	dir, err := os.Stat(htmlPath)
	if err != nil {
		return "", errors.New("html path does not exist")
	}
	if !dir.IsDir() {
		return "", errors.New("html path is not a directory")
	}
	indexHtml := filepath.Join(htmlPath, "index.html")
	htmlFile, err := os.Stat(indexHtml)
	if err != nil {
		return "", errors.New("html path does not have an index")
	}
	if htmlFile.IsDir() {
		return "", errors.New("html path index... is a directory? ")
	}
	return htmlPath, nil
}
