package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slices"

	"fmt"

)

var skipDirs = []string{
	"vendor",
	"node_modules",
	".next",
	".git",
}

// FileExists checks if the file exists
func FileExists(searchPath, name string) bool {
	return FindFile(searchPath, name) != ""
}

// FindFile searches for the file inside the path recursively
// and returns the full path of the file if found
// If multiple files exist, tries to return the one closest to root
func FindFile(searchPath, name string) string {
	files := FindAllFiles(searchPath, name)
	if len(files) == 0 {
		return ""
	}

	slices.SortFunc(files, func(a, b string) bool {
		return len(strings.Split(a, string(os.PathSeparator))) < len(strings.Split(b, string(os.PathSeparator)))
	})
	return files[0]
}

// FindAllFiles searches for the file inside the path recursively and returns all matches
func FindAllFiles(searchPath, name string) []string {
	found := make([]string, 0)
	_ = filepath.WalkDir(searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor directories
			if slices.Contains(skipDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() == name {
			found = append(found, p)
		}

		return nil
	})

	return found
}

// FindAllFiles searches for the file inside the path recursively and returns all matches
func FindAllFilesWithExtension(searchPath, extension string) []string {
	found := make([]string, 0)
	_ = filepath.WalkDir(searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor directories
			// if slices.Contains(skipDirs, d.Name()) {
			// 	return filepath.SkipDir
			// }
			// return nil
			return filepath.SkipDir
		}

		if strings.Contains(d.Name(), extension) {
			found = append(found, p)
		}

		// if d.Name() == name {
		// 	found = append(found, p)
		// }

		return nil
	})

	return found
}

func GetMapValue(keyPath []string, data map[string]interface{}) (value interface{}, ok bool) {
	if len(keyPath) == 0 {
		return data, true
	}

	for _, key := range keyPath[:len(keyPath)-1] {
		if value, ok = data[key]; !ok {
			return nil, false
		}

		if data, ok = value.(map[string]interface{}); !ok {
			return nil, false
		}
	}

	if value, ok = data[keyPath[len(keyPath)-1]]; !ok {
		return nil, false
	}

	return value, true
}

// GetJSONValue gets a value from a JSON file, by traversing the path given
func GetJSONValue(keyPath []string, filePath string, caseInsensitive bool) (value interface{}, ok bool) {
	fin, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer fin.Close()

	rawData, err := io.ReadAll(fin)
	if err != nil {
		return nil, false
	}

	if caseInsensitive {
		rawData = bytes.ToLower(rawData)
		for i := range keyPath {
			keyPath[i] = strings.ToLower(keyPath[i])
		}
	}

	var data map[string]interface{}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, false
	}

	return GetMapValue(keyPath, data)
}

// There are three common identifiers for finding FLASK_APP: the create_app or make_app factories, or the main Flask class.
func FindFlaskApp(workingDirectory string) (bool, string, error) {
	rootPyFiles := FindAllFilesWithExtension(workingDirectory,".py")
	for _, rootPyFilePath := range rootPyFiles {
		// create_app factory.
		f, err := os.Open(rootPyFilePath)
		if err != nil {
			return false, "", nil
		}
		defer f.Close()
		if ok, _, _ := ContainsStringInFile(f, "create_app(", false, false); ok {
			return true, filepath.Base(rootPyFilePath), nil
		} else {
			return false, "", nil
		}

		// run_app factory.
		// f, err := os.Open(rootPyFilePath)
		if err != nil {
			return false, "", nil
		}
		// defer f.Close()
		if ok, _, _ := ContainsStringInFile(f, "run_app(", false, false); ok {
			return true, filepath.Base(rootPyFilePath), nil
		} else {
			return false, "", nil
		}

		// Flask class.
		// f, err := os.Open(rootPyFilePath)
		if err != nil {
			return false, "", nil
		}
		// defer f.Close()
		if ok, _, _ := ContainsStringInFile(f, "= Flask(", false, false); ok {
			return true, filepath.Base(rootPyFilePath), nil
		} else {
			return false, "", nil
		}

	}

	return false, "", nil

}

// ContainsStringInFile checks if the given file contains the given string
func ContainsStringInFile(file io.Reader, target string, caseInsensitive bool, removeMatch bool) (bool, []string, error) {
	matchLines := []string{}
	stringInFile := false
	if caseInsensitive {
		target = strings.ToLower(target)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if caseInsensitive {
			if strings.Contains(strings.ToLower(scanner.Text()), target) {
				if removeMatch {
					matchLines = append(matchLines, strings.ReplaceAll(scanner.Text(), target, ""))
				} else {
					matchLines = append(matchLines)
				}
				stringInFile = true
			}
		} else {
			if strings.Contains(scanner.Text(), target) {
				if removeMatch {
					matchLines = append(matchLines, strings.ReplaceAll(scanner.Text(), target, ""))
				} else {
					matchLines = append(matchLines)
				}
				stringInFile = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return stringInFile, matchLines, err
	}

	return stringInFile, matchLines, nil
}

// GetTOMLValue gets a value from a TOML file, by traversing the path given
func GetTOMLValue(keyPath []string, filePath string, caseInsensitive bool) (value interface{}, ok bool) {
	fin, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer fin.Close()

	rawData, err := io.ReadAll(fin)
	if err != nil {
		return nil, false
	}

	if caseInsensitive {
		rawData = bytes.ToLower(rawData)
		for i := range keyPath {
			keyPath[i] = strings.ToLower(keyPath[i])
		}
	}

	var data map[string]interface{}
	err = toml.Unmarshal(rawData, &data)
	if err != nil {
		return nil, false
	}

	return GetMapValue(keyPath, data)
}

// DepInNestedRequirements follows the import path from a root requirements.txt file to up to 2 levels of imports
// 		to locate a particular dependency.
// 
// When using pip, a framework requirement can be nested across multiple files like so
// .
// ├── requirements
// │   ├── base.txt <- only this file contains Flask, for example
// │   ├── dev.txt  <- this file also imports base.txt
// │   └── prod.txt <- Flask _could_ be here, but this is also where this file could import base.txt
// └── requirements.txt <- this file imports production as an example
// 
// This is a common pattern, and without some logic like this we won't be able to detect many Python frameworks using pip + venv.
// 
// Note: this logic does still depend on a root requirements.txt file, as there's too much variability otherwise.
func DepInNestedRequirements(keyRequirement string, filePath string, caseInsensitive bool) (value interface{}, ok bool) {
	fin, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer fin.Close()

	// First check if Framework is imported in that root requirements.txt file.
	if ok, matchLines, _ := ContainsStringInFile(fin, keyRequirement, caseInsensitive, false); ok {
		return matchLines, true
	} else {
		f, err := os.Open(filePath)
		if err == nil {
			defer f.Close()

			// Then, check to see if there are imports (-r) in that root requirements.txt file.
			if ok, imports, _ := ContainsStringInFile(f, "-r ", true, true); ok {

				// If so, loop through all that are used.
				for _, importRequirementsFilePath := range imports {

					// Stash the import directory (i.e. requirements/), since everything is relative from here on.
					importDir := strings.Split(importRequirementsFilePath, "/")[0]
					f, err := os.Open(importRequirementsFilePath)
					if err == nil {
						defer f.Close()

						// Open the imported file, and check again for Framework.
						if ok, matchLines, _ := ContainsStringInFile(f, keyRequirement, caseInsensitive, false); ok {
							return matchLines, true

						} else {
							f, err := os.Open(importRequirementsFilePath)
							if err == nil {
								defer f.Close()

								// There may be one more layer of imports (i.e. requirements <- prod <- base ), so check.
								if ok, imports, _ := ContainsStringInFile(f, "-r ", true, true); ok {

									// If found, loop through all second-level imports.
									for _, importRequirementsFilePath := range imports {
										f, err := os.Open(fmt.Sprintf("%s/%s", importDir, importRequirementsFilePath))
										if err == nil {
											defer f.Close()

											// This is the deepest check, so check for the last time for Framework.
											if ok, matchLines, _ := ContainsStringInFile(f, keyRequirement, caseInsensitive, false); ok {
												return matchLines, true
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, false
}
