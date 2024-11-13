package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	colorReset          = "\033[0m"
	colorRed            = "\033[31m"
	colorGreen          = "\033[32m"
	colorYellow         = "\033[33m"
	errResult           = "error"
	infoResult          = "info"
	warnResult          = "warn"
	defaultNMOutputFile = "nm_output.txt"

	invalidProfileExitCode = 4
)

// Rule represents a single rule for checking the binary.
type Rule struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	Regex          string `yaml:"regex"`
	FoundResult    string `yaml:"found_result"`
	NotFoundResult string `yaml:"not_found_result"`
}

// Rules contains all the rules for checking the binary.
type Rules struct {
	NmRules        []Rule `yaml:"nm_rules"`
	CallGraphRules []Rule `yaml:"call_graph_rules"`
}

// RulesConfig represents the configuration for all rules.
type RulesConfig struct {
	Rules Rules `yaml:"rules"`
}

// GetResultColor returns the color code for a given result severity.
func GetResultColor(result string) string {
	switch result {
	case errResult:
		return colorRed
	case warnResult:
		return colorYellow
	case infoResult:
		return colorGreen
	default:
		return colorReset
	}
}

// PrintMatch prints the result of a rule match.
func PrintMatch(matchResult bool, rule Rule) {

	var severity string
	var message string

	if matchResult {
		severity = rule.FoundResult
	} else {
		severity = rule.NotFoundResult
	}
	if severity != infoResult {
		message = rule.Description
	}

	color := GetResultColor(severity)
	var foundMessage string
	if matchResult {
		foundMessage = "found"
	} else {
		foundMessage = "not found"
	}

	// Print check results
	if severity == errResult {
		fmt.Printf("Check: %s %s - %s.  %s %s. \nTo fix: %s %s\n", color, rule.Name, severity, rule.Name, foundMessage, message, colorReset)
	} else {
		fmt.Printf("Check: %s %s - %s.  %s %s. %s\n", color, rule.Name, severity, rule.Name, foundMessage, colorReset)

	}

}

// GenerateNMFile generates the nm output file for the given binary file path.
func GenerateNMFile(binaryFilePath string) (string, error) {
	// execute go tool nm on the binary file with profiles
	cmd := exec.Command("go", "tool", "nm", binaryFilePath)
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing nm command: %v", err)
	}
	err = os.WriteFile(defaultNMOutputFile, stdout, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing nm output to file: %v", err)
	}

	nmFileOutput, error := os.ReadFile(defaultNMOutputFile)
	if error != nil {
		return "", fmt.Errorf("error reading nm output file: %v", error)
	}
	nmFileOutputString := string(nmFileOutput)
	return nmFileOutputString, nil
}

// LoadRulesFromFile loads the rules from the specified YAML file.
func LoadRulesFromFile(filename string) (*RulesConfig, error) {
	// Load the pattern file
	if _, err := os.Stat(filename); err != nil {
		log.Printf("Pattern file is missing: %v\n", err)
		os.Exit(invalidProfileExitCode)
	}

	// Read the YAML file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// Parse the YAML file
	var config RulesConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Failed to parse YAML file: %v", err)
		return nil, err
	}
	return &config, nil
}

// GetCompiledRegexs compiles the regex patterns for the given rules.
func GetCompiledRegexs(rules []Rule) map[string]*regexp.Regexp {
	regexes := make(map[string]*regexp.Regexp)
	for _, rule := range rules {
		regexes[rule.Name] = regexp.MustCompile(rule.Regex)
	}
	return regexes
}
