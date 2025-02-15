package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Define flags
	rootDir := flag.String("dir", "", "Directory to monitor") // Will default to the current working directory if no flag is defined
	generateBaseline := flag.String("baseline", "", "Generates the baseline json file")
	compareToBaselineFile := flag.String("compare", "", "Takes a file name. Will compare whatever that baseline file took as root and re-scan against the baseline.")

	flag.Parse()

	// Check to make sure that both -dir and -baseline are used together
	if *rootDir != "" && *generateBaseline == "" {
		log.Fatal("ERROR: if you are defining a directory, you must also generate a baseline file with -baseline <filename>. Otherwise the program won't do anything")
	}
	if *generateBaseline != "" && *rootDir == "" {
		log.Fatal("ERROR: if you are generating a baseline file, you must define the root directory with -dir <directory path>")
	}

	// jsonData to be used when generating a baseline .json file
	var jsonData []byte

	if *rootDir != "" {
		fmt.Println("[+] Getting file paths and hash values")
		// Check for errors if the root directory to be scanned doesn't exist
		if _, err := os.Stat(*rootDir); os.IsNotExist(err) {
			log.Fatal("directory does not exist:", err)
		} else if err != nil {
			log.Fatal("error checking directory:", err)
		}

		// Get paths to all files in the root directory and any subdirectories
		filepaths, err := getFilePaths(*rootDir)
		if err != nil {
			log.Fatal("error getting file paths:", err)
		}

		// Generate baseline files struct to be written to json when creating baseline
		filesHashed, err := getHashValues(*rootDir, filepaths)
		if err != nil {
			log.Fatal("error generating a hash map:", err)
		}

		// Write to a json
		jsonData, err = writeToJson(filesHashed)
		if err != nil {
			log.Fatal("error writing struct to json:", err)
		}
	}

	// If the generate baseline flag is used, write json data to a .json file
	if *generateBaseline != "" {
		fmt.Println("[+] Generating baseline file...")
		returnMessage, err := generateBaselineFunc(*generateBaseline, jsonData)
		if err != nil {
			log.Fatal("error generating baseline file:", err)
		}
		fmt.Println(returnMessage)
	}

	// If compare to baseline flag is used, take the root defined in the baseline file and re-scan for changes, additions, or removed files
	if *compareToBaselineFile != "" {
		fmt.Println("[+] Comparing current files with baseline...")
		baseLineFilePath, err := getComparePath(*compareToBaselineFile)
		if err != nil {
			log.Fatal("error reading baseline file. Make sure you're file name is correct:", err)
		}
		// Write differences to Difference struct
		differences, err := compareBaselineFunc(baseLineFilePath)
		if err != nil {
			log.Fatal("error comparing with baseline:", err)
		}

		// Take a Difference struct and format/print it
		printDifferences(differences)

		// Write differences to a new json file
		err = writeDifToJson(*compareToBaselineFile, baseLineFilePath.BaselinePath, differences)
		if err != nil {
			log.Fatal("error writing differences to json file:", err)
		}
	}
}
