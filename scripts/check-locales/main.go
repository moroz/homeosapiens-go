package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"strings"
)

func resolveGitDir() (string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := start

	for {
		gitDir := filepath.Join(dir, ".git")
		stat, _ := os.Stat(gitDir)
		if stat != nil && stat.IsDir() {
			return dir, nil
		}
		if filepath.Dir(dir) == dir {
			return "", fmt.Errorf("Git directory not found in parent directories of %s", start)
		}
		dir = filepath.Dir(dir)
	}
}

func listLocaleFiles(parentDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(parentDir, "*.json"))
}

func resolveKey(prefix []string, key string) string {
	return strings.Join(appendKey(prefix, key), ".")
}

func appendKey(prefix []string, key string) []string {
	segments := append([]string{}, prefix...)
	segments = append(segments, key)
	return segments
}

func mergeMap(m map[string]any, prefix []string) map[string]string {
	result := make(map[string]string)
	for key, value := range m {
		switch value := value.(type) {
		case string:
			resolvedKey := resolveKey(prefix, key)
			result[resolvedKey] = value

		case map[string]any:
			newPrefix := appendKey(prefix, key)
			merged := mergeMap(value, newPrefix)
			maps.Copy(result, merged)
		}
	}
	return result
}

func decodeFile(path string) (map[string]any, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result map[string]any
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	gitDir, err := resolveGitDir()
	if err != nil {
		log.Fatal(err)
	}

	parentDir := filepath.Join(gitDir, "i18n")
	localeFiles, err := listLocaleFiles(parentDir)
	if err != nil {
		log.Fatal(err)
	}

	allLocales := make(map[string]map[string]string)

	for _, file := range localeFiles {
		basepath := filepath.Base(file)
		locale := strings.TrimSuffix(basepath, filepath.Ext(basepath))

		data, err := decodeFile(file)
		if err != nil {
			log.Fatal(err)
		}

		allLocales[locale] = mergeMap(data, nil)

		fmt.Printf("%+v\n", allLocales[locale])
	}

}
