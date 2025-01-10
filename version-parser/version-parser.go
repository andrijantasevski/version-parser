package lastversionparser

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	pathFlag    string
	versionFlag string
	direction   string
	help string
)

func init() {
	flag.StringVar(&pathFlag, "path", "", "With the --path parameter, we can define the path to the file for parsing.\nIf you are using Windows, make sure to provide double backslashes (\\\\) or a forward slash (/).")
	flag.StringVar(&versionFlag, "version-flag", "", "With the --version-flag parameter, we can define what flag we want to parse for.\nFor example, if we want to parse the version number from the line 'asfafasigTX10', we can use the flag 'TX'.")
	flag.StringVar(&direction, "direction", "", "With the direction parameter, we can define the direction of the parsing. We can choose between 'min or 'max'.")

	flag.Parse()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsPathValid(path string) (bool, bool) {
	info, err := os.Stat(path)

	if err == nil {
		return true, info.IsDir()
	}

	return false, false
}

func IsFileExtensionValid(path string) bool {
	ext := filepath.Ext(path)

	return ext == ".txt"
}

func BuildVersionRegex(versionFlag string) (*regexp.Regexp, error) {
	if versionFlag == "" {
		return nil, fmt.Errorf("versionFlag cannot be empty")
	}
	pattern := fmt.Sprintf(`%s\d+`, regexp.QuoteMeta(versionFlag))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %v", err)
	}
	return re, nil
}

func FindLastVersionInFile(path, versionFlag string) (int, error) {
	re, err := BuildVersionRegex(versionFlag)

	if err != nil {
		return 0, fmt.Errorf("error building regex: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	uniqueMatches := make(map[string]bool)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		matches := re.FindAllString(line, -1)

		for _, match := range matches {
			uniqueMatches[match] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %v", err)
	}

	if len(uniqueMatches) == 0 {
		return 0, fmt.Errorf("no matches found for the version flag %q", versionFlag)
	}

	result := make([]int, 0, len(uniqueMatches))

	for match := range uniqueMatches {
		fmt.Println(match)
		versionStr := strings.TrimPrefix(match, versionFlag)
		version, err := strconv.Atoi(versionStr)

		if err != nil {
			return 0, fmt.Errorf("invalid version number in match %q: %v", match, err)
		}

		result = append(result, version)
	}

	if direction == "max" {
		return slices.Max(result), nil
	}

	return slices.Min(result), nil
}

func ParseForVersion() {
	if help != "" {
		flag.Usage()
		os.Exit(0)
	}

	isPathValid, isPathDir := IsPathValid(pathFlag)

	if !isPathValid {
		fmt.Println("\nThe path you have provided is not valid. Please provide a valid path.\nIf you are using Windows, make sure to provide double backslashes (\\\\) or a forward slash (/).")
		os.Exit(1)
	}

	if isPathValid && isPathDir {
		fmt.Println("\nThe path you have provided is a directory. Please provide a valid path to a .txt file.")
		os.Exit(1)
	}

	if !IsFileExtensionValid(pathFlag) {
		fmt.Println("\nThe file extension you have provided is not valid. Please provide a valid file with a .txt extension.")
		os.Exit(1)
	}

	if versionFlag == "" {
		fmt.Println("\nPlease provide a valid version flag.")
		os.Exit(1)
	}

	direction = strings.ToLower(direction)
	if direction != "max" && direction != "min" {
		fmt.Println("\n Please provice a valid direction. You can choose between 'min' or 'max'.")
		os.Exit(1)
	}

	fmt.Printf("\nParsing all instances in the .txt file with the %q version flag.", versionFlag)

	lastVersion, err := FindLastVersionInFile(pathFlag, versionFlag)

	if err != nil {
		fmt.Println("\nError:", err)
		os.Exit(1)
	}

	fmt.Println("\nThe last version is:", lastVersion)
}
