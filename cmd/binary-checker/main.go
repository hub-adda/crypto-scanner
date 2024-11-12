package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	utils "crypto-scanner/internal/pkg/utils"
)

const (
	defaultProfileFilePath = "profiles/default.yaml"
	invalidArgsExitCode    = iota + 1
	invalidBinaryExitCode
	invalidProfileExitCode
	unsafeCryptoExitCode
	toolMissingExitCode
)

type RuntimeConfig struct {
	BinaryFilePath  string
	ProfileFilePath string
}

func main() {
	// Ensure at least one argument is provided
	if len(os.Args) < 1 {
		printUsage()
		os.Exit(invalidArgsExitCode)
	}

	// Check if Go compiler is installed
	if !commandExists("go") {
		fmt.Println("Go compiler could not be found. Please install Go to proceed.")
		os.Exit(invalidArgsExitCode)
	}

	// Parse command line arguments
	runtimeConfig := parseCommandLineArgs()

	// Print the profile file name and scanned binary file
	fmt.Printf("Using profile file: %s\n", runtimeConfig.ProfileFilePath)
	fmt.Printf("Scanning binary file: %s\n", runtimeConfig.BinaryFilePath)

	// Load rules from the profile file
	rules, err := utils.LoadRulesFromFile(runtimeConfig.ProfileFilePath)
	if err != nil {
		fmt.Printf("Error loading pattern file: %v\n", err)
		os.Exit(invalidProfileExitCode)
	}

	// Check if the binary file is a valid Go binary
	err = checkGolangLinuxBinary(runtimeConfig.BinaryFilePath)
	if err != nil {
		fmt.Printf("Error checking binary file: %v\n", err)
		os.Exit(invalidBinaryExitCode)
	}

	// Check the binary file with nm tool patterns
	checkNMRules(rules.Rules.NmRules, runtimeConfig)
}

// printUsage prints the usage help message
func printUsage() {
	fmt.Printf("Usage: %s -binary <binary_file> [-profile <profile_file>]\n", os.Args[0])
	fmt.Printf("Default profile file: %s\n", defaultProfileFilePath)
}

// checkNMRules checks the binary file against nm tool patterns
func checkNMRules(rules []utils.Rule, config RuntimeConfig) {
	output, err := utils.GenerateNMFile(config.BinaryFilePath)
	if err != nil {
		fmt.Printf("Error generating nm file: %v\n", err)
		os.Exit(toolMissingExitCode)
	}

	// Compile Regex from rules
	regexes := utils.GetCompiledRegexs(rules)

	for _, rule := range rules {
		reg := regexes[rule.Name]
		matchFound := reg.MatchString(output)
		utils.PrintMatch(matchFound, rule)
	}
}

// parseCommandLineArgs parses the command line arguments and returns a RuntimeConfig object
func parseCommandLineArgs() RuntimeConfig {
	config := RuntimeConfig{}
	flag.StringVar(&config.ProfileFilePath, "profile", defaultProfileFilePath, "Path to the profile file containing the rules.")
	flag.StringVar(&config.BinaryFilePath, "binary", "", "Path to the binary file to be checked.")
	flag.Parse()
	return config
}

// checkGolangLinuxBinary checks if the file is a valid Go binary
func checkGolangLinuxBinary(filePath string) error {
	// Check if the file exists
	rule := utils.Rule{
		Name:           "The file exists",
		Description:    "Check if the file exists",
		Regex:          "",
		FoundResult:    "info",
		NotFoundResult: "error",
	}
	if _, err := os.Stat(filePath); err != nil {
		utils.PrintMatch(false, rule)
		return fmt.Errorf("binary file not found: %v", err)
	}

	// Check if the file is a valid Go binary
	rule = utils.Rule{
		Name:           "The file is a valid Go binary",
		Description:    "Check if the file is a valid Go binary",
		Regex:          "",
		FoundResult:    "info",
		NotFoundResult: "error",
	}
	cmd := exec.Command("go", "tool", "nm", filePath)
	stdout, err := cmd.Output()
	if err != nil {
		utils.PrintMatch(false, rule)
		return fmt.Errorf("file is invalid Go binary: %v (%s)", err, stdout)
	}
	utils.PrintMatch(true, rule)

	return nil
}

// commandExists checks if a command exists in the system
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
