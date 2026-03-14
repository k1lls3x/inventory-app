package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func loadDotEnv() {
	for _, path := range envCandidatePaths() {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}

		if err := godotenv.Load(path); err == nil {
			log.Printf("Загружен конфиг из %s", path)
			return
		}
	}
}

func envCandidatePaths() []string {
	var candidates []string

	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, ".env")
		candidates = append(candidates, ascendEnvPaths(cwd, 4)...)
	}

	if executablePath, err := os.Executable(); err == nil {
		candidates = append(candidates, ascendEnvPaths(filepath.Dir(executablePath), 4)...)
	}

	return uniqueStrings(candidates)
}

func ascendEnvPaths(startDir string, limit int) []string {
	var paths []string
	currentDir := startDir

	for i := 0; i < limit; i++ {
		paths = append(paths, filepath.Join(currentDir, ".env"))
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return paths
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}
