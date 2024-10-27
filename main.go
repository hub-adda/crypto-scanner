package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Pattern struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	Regex          string `yaml:"regex"`
	FoundResult    string `yaml:"found_result"`
	NotFoundResult string `yaml:"not_found_result"`
}

type PatternConfig struct {
	Patterns []Pattern `yaml:"patterns"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <nm_output_file> [pattern_file]")
		return
	}

	nmOutputFile := os.Args[1]
	patternFile := "default.yaml"
	if len(os.Args) > 2 {
		patternFile = os.Args[2]
	}

	patterns, err := loadPatterns(patternFile)
	if err != nil {
		fmt.Printf("Error loading pattern file: %v\n", err)
		return
	}

	file, err := os.Open(nmOutputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Compile the patterns into regex objects
	regexes := make(map[string]*regexp.Regexp)
	for _, pattern := range patterns {
		regexes[pattern.Name] = regexp.MustCompile(pattern.Regex)
	}

	// Initialize maps to track pattern matches
	foundPatterns := make(map[string]bool)
	for _, pattern := range patterns {
		foundPatterns[pattern.Name] = false
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for name, regex := range regexes {
			if regex.MatchString(line) {
				foundPatterns[name] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Print results
	fmt.Println("Pattern scan results:")
	for _, pattern := range patterns {
		if foundPatterns[pattern.Name] {
			fmt.Printf("Pattern '%s' found: %s\n", pattern.Name, pattern.FoundResult)
		} else {
			fmt.Printf("Pattern '%s' not found: %s\n", pattern.Name, pattern.NotFoundResult)
		}
	}
}

func loadPatterns(filename string) ([]Pattern, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config PatternConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config.Patterns, nil
}
