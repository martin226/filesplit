package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	flag "github.com/spf13/pflag"
)

func parseMode() string {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: filesplit [mode] [-F, --file] [-N, --number] \n")
		fmt.Println("Available modes:")
		fmt.Println("  file   Split the file into a specific number of sub-files.")
		fmt.Println("  line   Split the file into sub-files with a specific number of lines per file.")
		os.Exit(1)
	}

	mode := os.Args[1]
	
	if mode == "file" || mode == "line" {
		return mode
	} else {
		fmt.Printf("%v: unknown mode\n", os.Args[1])
		os.Exit(1)
		return ""
	}

}

func parseArguments() (*int, *string) {
	numberPtr := flag.IntP("number","N", 0, "Number of output files. Must be greater than 0. (Required)")
	fPathPtr := flag.StringP("file", "F", "", "Path to file. (Required)")
	flag.Parse()

	if *fPathPtr == "" || *numberPtr < 1 {
		fmt.Printf("Usage: filesplit [mode] [-F, --file] [-N, --number] \n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	return numberPtr, fPathPtr
}

func getLines(scanner *bufio.Scanner) []string {
	result := []string{}

	for scanner.Scan(){
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

func writeLines(lines []string, fName string) error {
	f, err := os.Create(fName) 
	if err != nil {
		return err
	}
	defer f.Close()
	for _, line := range lines {
		fmt.Fprintln(f, line)
	}
	return nil
}

func mkOutputDir(fName string) (string, error) {
	dName := fmt.Sprintf("%v_%v", fName, time.Now().Unix())

	err := os.Mkdir(dName, 0755)
	if err != nil {
		return "", err
	}
	return dName, nil
}

func mkOutputFiles(lineCount int, linesPerFile int, lines []string, dName string, fileExt string) error {
	var j int
	var fileN int
	for i := 0; i < lineCount; i += linesPerFile{
		fileN++
		j += linesPerFile
		if j > lineCount {
			j = lineCount
		}
		chunk := lines[i:j]
		err := writeLines(chunk, fmt.Sprintf("%v/output_%v.%v", dName, fileN, fileExt))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: filesplit [mode] [-F, --file] [-N, --number] \n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parse mode
	mode := parseMode()

	// Parse arguments
	numberPtr, fPathPtr := parseArguments()

	// Read file
	f, err := os.Open(*fPathPtr)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	lines := getLines(scanner)

	lineCount := len(lines)
	if lineCount == 0 {
		return
	}

	var linesPerFile int
	if mode == "file" {
		linesPerFile = (lineCount + *numberPtr - 1) / *numberPtr
	} else {
		linesPerFile = *numberPtr
	}

	// Create output directory
	dName, err2 := mkOutputDir(filepath.Base(*fPathPtr))
	if err2 != nil {
		panic(err2)
	}

	// Write output files
	err3 := mkOutputFiles(lineCount, linesPerFile, lines, dName, filepath.Ext(*fPathPtr))
	if err3 != nil {
		panic(err3)
	}
}