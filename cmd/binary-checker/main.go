package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

// main is the entry point of the application.
func main() {
	printLogo()

	if len(os.Args) < 1 {
		printUsage()
		os.Exit(invalidArgsExitCode)
	}

	if !commandExists("go") {
		fmt.Println("Go compiler could not be found. Please install Go to proceed.")
		os.Exit(invalidArgsExitCode)
	}

	runtimeConfig := parseCommandLineArgs()

	if _, err := os.Stat(runtimeConfig.BinaryFilePath); os.IsNotExist(err) {
		fmt.Printf("Binary file does not exist: %s\n", runtimeConfig.BinaryFilePath)
		printUsage()
		os.Exit(invalidBinaryExitCode)
	}

	if _, err := os.Stat(runtimeConfig.ProfileFilePath); os.IsNotExist(err) {
		fmt.Printf("Profile file does not exist: %s\n", runtimeConfig.ProfileFilePath)
		printUsage()
		os.Exit(invalidProfileExitCode)
	}

	fmt.Printf("Using profile file: %s\n", runtimeConfig.ProfileFilePath)
	fmt.Printf("Scanning binary file: %s\n", runtimeConfig.BinaryFilePath)

	rules, err := utils.LoadRulesFromFile(runtimeConfig.ProfileFilePath)
	if err != nil {
		fmt.Printf("Error loading pattern file: %v\n", err)
		os.Exit(invalidProfileExitCode)
	}

	fmt.Printf("\nChecking Binary\n")
	err = checkGolangLinuxBinary(runtimeConfig.BinaryFilePath)
	if err != nil {
		fmt.Printf("Error checking binary file: %v\n", err)
		os.Exit(invalidBinaryExitCode)
	}

	checkNMRules(rules.Rules.NmRules, runtimeConfig)
}

// checkNMRules checks the binary file against the provided nm rules.
func checkNMRules(rules []utils.Rule, config RuntimeConfig) {
	output, err := utils.GenerateNMFile(config.BinaryFilePath)
	if err != nil {
		fmt.Printf("Error generating nm file: %v\n", err)
		os.Exit(toolMissingExitCode)
	}

	regexes := utils.GetCompiledRegexs(rules)

	for _, rule := range rules {
		reg := regexes[rule.Name]
		matchFound := reg.MatchString(output)
		utils.PrintMatch(matchFound, rule)
	}
}

// parseCommandLineArgs parses the command line arguments and returns a RuntimeConfig.
func parseCommandLineArgs() RuntimeConfig {
	config := RuntimeConfig{}
	flag.StringVar(&config.ProfileFilePath, "profile", "", "Path to the profile file containing the rules.")
	flag.StringVar(&config.BinaryFilePath, "binary", "", "Path to the binary file to be checked.")
	flag.Parse()

	config.ProfileFilePath = convertToAbsolutePath(config.ProfileFilePath, invalidProfileExitCode)
	config.BinaryFilePath = convertToAbsolutePath(config.BinaryFilePath, invalidBinaryExitCode)

	return config
}

// convertToAbsolutePath converts a relative file path to an absolute path.
func convertToAbsolutePath(path string, exitCode int) string {
	if !filepath.IsAbs(path) && path != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("Error converting path to absolute: %v\n", err)
			os.Exit(exitCode)
		}
		return absPath
	}
	return path
}

// checkGolangLinuxBinary checks if the provided file is a valid Go binary.
func checkGolangLinuxBinary(filePath string) error {
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

// commandExists checks if a command exists in the system's PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// printUsage prints the usage instructions for the application.
func printUsage() {
	fmt.Printf("Usage: %s -binary <binary_file> [-profile <profile_file>]\n", os.Args[0])
	fmt.Printf("Default profile file: %s\n", defaultProfileFilePath)
}

// printLogo prints the ASCII art logo of the application.
func printLogo() {
	fmt.Println("+--------------------------+")
	fmt.Println("| Binary Safe Checker v0.1 |")
	fmt.Println("+--------------------------+")
	fmt.Printf("Tool to validate the cryptogrpahic functionallity used in a Go binary \n\n")
}
