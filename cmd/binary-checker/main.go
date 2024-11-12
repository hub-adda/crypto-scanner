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
)

type RuntimeConfig struct {
	BinaryFilePath  string
	ProfileFilePath string
}

func main() {

	if len(os.Args) < 1 {

		fmt.Println("Usage: go run main.go <binary_file> [pattern_file]")

		os.Exit(invalidArgsExitCode)
	}

	// Check if go compiler is installed
	if !commandExists("go") {
		fmt.Println("go compiler could not be found. Please install Go to proceed.")
		os.Exit(invalidArgsExitCode)
	}

	// Parse command line arguments

	runtimeConfig := parseCommandLineArgs()

	// Print the profile file name and scanned binary file
	fmt.Printf("Using profile file: %s\n", runtimeConfig.ProfileFilePath)
	fmt.Printf("Scanning binary file: %s\n", runtimeConfig.BinaryFilePath)

	rules, err := utils.LoadRulesFromFile(runtimeConfig.ProfileFilePath)

	if err != nil {
		fmt.Printf("Error loading pattern file: %v\n", err)
		os.Exit(4)
	}

	// Check if the binary file is a valid go binary
	err = checkGolangLinuxBinary(runtimeConfig.BinaryFilePath)
	if err != nil {
		fmt.Printf("Error checking binary file: %v\n", err)
		os.Exit(invalidBinaryExitCode)
	}

	// Check the binary file with nm tool patterns that are allowed or
	checkNMRules(rules.Rules.NmRules, runtimeConfig)
}

func checkNMRules(rules []utils.Rule, config RuntimeConfig) {
	output, err := utils.GenerateNMFile(config.BinaryFilePath)
	if err != nil {
		fmt.Printf("Error generating nm file: %v\n", err)
		os.Exit(3)
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
	flag.StringVar(&config.ProfileFilePath, "profile", "profiles/default.yaml", "Path to the profile file containing the rules.")
	flag.StringVar(&config.BinaryFilePath, "binary", "", "Path to the binary file to be checked.")
	flag.Parse()
	return config
}

// checkGolangLinuxBinary checks if the file is a valid go binary
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
	// Check if the file is a valid go binary
	rule = utils.Rule{
		Name:           "The file is a valid go binary",
		Description:    "Check if the file is a valid go binary",
		Regex:          "",
		FoundResult:    "info",
		NotFoundResult: "error",
	}

	cmd := exec.Command("go", "tool", "nm", filePath)
	stdout, err := cmd.Output()
	if err != nil {
		utils.PrintMatch(false, rule)
		return fmt.Errorf("file is invalid go binary: %v (%s)", err, stdout)
	}
	utils.PrintMatch(true, rule)

	// validate the binary is built for linux
	// rule = utils.Rule{
	// 	Name:           "The file is built for Linux OS",
	// 	Description:    "Check if the file is built for Linux OS",
	// 	Regex:          "",
	// 	FoundResult:    "info",
	// 	NotFoundResult: "error",
	// }
	// isElf := isELFLinux(filePath)
	// utils.PrintMatch(isElf, rule)
	return nil
}

// commandExists checks if a command exists in the system
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ******************************************
// to be replaced with go version tool checks
// ******************************************

// func isELFLinux(filePath string) bool {

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return false
// 	}
// 	defer file.Close()

// 	magicNumber := make([]byte, 4)
// 	_, err = file.Read(magicNumber)
// 	if err != nil {
// 		return false
// 	}

// 	// Go binaries have a specific magic number
// 	expectedMagicNumber := []byte{0x7f, 'E', 'L', 'F'} // Example for ELF binaries

// 	return string(magicNumber) == string(expectedMagicNumber)
// }
