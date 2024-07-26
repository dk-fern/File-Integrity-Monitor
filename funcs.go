package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Get all file paths in the given root directory. Only returns full file paths and no directories
func getFilePaths(rootDir string) ([]string, error) {
	var paths []string

	// Walk through the directory tree
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If the file doesn't exist or there is a permission error, skip the file
			if os.IsNotExist(err) || os.IsPermission(err) {
				// fmt.Printf("Skipping path %v: %v\n", path, err) << Uncomment if desire is to see what is being skipped
				return nil
			}
			// For other errors, print and continue
			fmt.Printf("Error accessing the path %v: %v\n", path, err)
			return nil
		}

		if info.IsDir() {
			// Print each file or directory's path
			return nil
		}

		if info.Mode()&os.ModeSocket != 0 {
			// fmt.Printf("Skipping socket: %v\n", path) << Uncomment if desire is to see what is being skipped
			return nil
		}

		if info.Mode()&os.ModeNamedPipe != 0 {
			// fmt.Printf("Skipping named pipe: %v\n", path) << Uncomment if desire is to see what is being skipped
			return nil
		}

		if info.Mode()&os.ModeDevice != 0 {
			// fmt.Printf("Skipping device file: %v\n", path) << Uncomment if desire is to see what is being skipped
			return nil
		}

		paths = append(paths, path)
		return nil

	})
	if err != nil {
		fmt.Println("Error walking the path:", err)
	}

	return paths, nil
}

// Get hash values of all files in a given path list. Takes the initial root directory as that is used for further identification when comparing files
func getHashValues(rootDir string, pathList []string) (BaselineFileList, error) {
	var hashMap []File

	for _, path := range pathList {
		// Check if the path is a directory
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Log files that raise a "do not exist" error
				// fmt.Printf("File %v does not exist, skipping...\n", path) << Uncomment if desire is to see what is being skipped
				continue
			}
			// Log other errors
			fmt.Printf("Error accessing file %v: %v\n", path, err)
			continue
		}

		// Ocassionally a directory will still make it through the last function, so this will check again.
		if info.IsDir() {
			// fmt.Printf("Path %v is a directory, skipping...\n", path) << Uncomment if desire is to see what is being skipped
			continue
		}

		// Process the file
		hasher := sha256.New()
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %v: %v\n", path, err)
			continue
		}
		defer file.Close()

		if _, err := io.Copy(hasher, file); err != nil {
			fmt.Printf("Error reading file %v: %v\n", path, err)
			continue
		}

		fileHash := hex.EncodeToString(hasher.Sum(nil))
		// Append file struct to the slice of File struct
		hashMap = append(hashMap, File{Path: path, Hash: fileHash})
	}

	fmt.Println("[+] Finished processing files...")

	// Create BaselineFileList struct and return it
	baselineFileList := BaselineFileList{
		BaselinePath: rootDir,
		Files:        hashMap,
	}
	return baselineFileList, nil
}

// Takes the baselineFileList struct and returns a json object
func writeToJson(baselineFileList BaselineFileList) ([]byte, error) {
	jsonFile, err := json.MarshalIndent(baselineFileList, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("error creating json: %v", err)
	}
	return jsonFile, nil
}

// Creates a new .json file and returns the file name
func generateBaselineFunc(definingFilename string, jsonData []byte) (string, error) {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	filename := fmt.Sprintf("%s_Baseline: %s.json", definingFilename, currentDate)

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("error creating new baseline file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(string(jsonData))
	if err != nil {
		return "", fmt.Errorf("error writing json to file: %v", err)
	}

	return fmt.Sprintln("[+] Successfully generated baseline file"), nil
}

// Takes a baseline file, and marshals it to a BaselineFileList struct to then compare with a new scan.
func getComparePath(filename string) (BaselineFileList, error) {
	// Open the json file
	jsonFile, err := os.Open(filename)
	if err != nil {
		return BaselineFileList{}, nil
	}
	defer jsonFile.Close()

	// read the json file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return BaselineFileList{}, nil
	}

	// create instance of a BaselineFileList and marshal the json data to the object
	var baselineFileList BaselineFileList
	err = json.Unmarshal(byteValue, &baselineFileList)
	if err != nil {
		return BaselineFileList{}, nil
	}

	return baselineFileList, nil

}

// Take a baseline struct, re-scan the same root path and return a Difference struct which shows changed, added, or removed files
func compareBaselineFunc(baseLine BaselineFileList) (Difference, error) {
	// Get baseline path and file list
	baselineRootPath := baseLine.BaselinePath
	baselineFiles := baseLine.Files

	// Create a map of baseline file paths to their hash values
	baselineMap := make(map[string]string)
	for _, file := range baselineFiles {
		baselineMap[file.Path] = file.Hash
	}

	// Get new file paths based on the baseline root path
	comparePaths, err := getFilePaths(baselineRootPath)
	if err != nil {
		return Difference{}, err
	}
	// Get new hash values of the files based on the new file paths
	compareFiles, err := getHashValues(baselineRootPath, comparePaths)
	if err != nil {
		return Difference{}, err
	}
	// Create a map of compare file paths and their hash values
	compareMap := make(map[string]string)
	for _, file := range compareFiles.Files {
		compareMap[file.Path] = file.Hash
	}

	// Create instances of Differences elements
	var hashDifferencesSlice []string
	var addedFilesSlice []string
	var removedFilesSlice []string

	// Identify differences in hash values and added files
	for path, hash := range compareMap {
		if baselineHash, exists := baselineMap[path]; exists {
			if hash != baselineHash {
				hashDifferencesSlice = append(hashDifferencesSlice, path)
			}
		} else {
			addedFilesSlice = append(addedFilesSlice, path)
		}
	}

	// Identify removed files
	for path := range baselineMap {
		if _, exists := compareMap[path]; !exists {
			removedFilesSlice = append(removedFilesSlice, path)
		}
	}

	differences := Difference{
		HashDifferences: hashDifferencesSlice,
		AddedFiles:      addedFilesSlice,
		RemovedFiles:    removedFilesSlice,
	}

	return differences, nil
}

// Use function to format and print Difference
func printDifferences(dif Difference) {
	fmt.Println("[+] Results:")
	fmt.Println("\n~~~Hash Differences~~~")
	if dif.HashDifferences == nil {
		fmt.Println("No hash differences found")
	}
	for n, hashDif := range dif.HashDifferences {
		fmt.Println(n+1, "-", hashDif)
	}

	fmt.Println("~~~Files Added~~~")
	if dif.AddedFiles == nil {
		fmt.Println("No files added")
	}
	for n, filesAdded := range dif.AddedFiles {
		fmt.Println(n+1, "-", filesAdded)
	}

	fmt.Println("~~~Files Removed~~~")
	if dif.RemovedFiles == nil {
		fmt.Println("No files removed")
	}
	for n, filesRemoved := range dif.RemovedFiles {
		fmt.Println(n+1, "-", filesRemoved)
	}
}

// Write differences to a new json file
func writeDifToJson(baselineFilename string, rootPath string, dif Difference) error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")
	compareFilename := fmt.Sprintf("Compare %v: %v", currentDate, baselineFilename)

	compareFileList := CompareFileList{
		BaselinePath: rootPath,
		Differences:  dif,
	}

	// Write dif to json
	jsonData, err := json.MarshalIndent(compareFileList, "", "    ")
	if err != nil {
		return err
	}

	file, err := os.Create(compareFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(string(jsonData))
	if err != nil {
		return err
	}

	return nil
}
