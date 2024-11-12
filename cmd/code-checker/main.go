package main

import (
	"crypto-scanner/internal/pkg/utils"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const (
	defaultProfileFilePath = "profiles/default.yaml"
	defaultNMOutputFile    = "nm_output.txt"
	invalidArgsExitCode    = iota + 1
	invalidBinaryExitCode
	invalidProfileExitCode
	unsafeCryptoExitCode
)

type RuntimeConfig struct {
	SrcFilePath     string
	ProfileFilePath string
}

func main() {

	if len(os.Args) < 1 {
		fmt.Println("Usage: go run main.go <main_source_file> [pattern_file]")
		os.Exit(invalidArgsExitCode)
	}

	// Parse command line arguments
	runtimeConfig := parseCommandLineArgs()

	log.Printf("Scanning file: %s for unsafe implementations\nwith profile:%s\n", runtimeConfig.SrcFilePath, runtimeConfig.ProfileFilePath)

	config, err := utils.LoadRulesFromFile(runtimeConfig.ProfileFilePath)
	if err != nil {
		fmt.Printf("Error loading pattern file: %v\n", err)
		os.Exit(4)
	}

	// Scan the code call graph
	filePath := runtimeConfig.SrcFilePath

	commentFipsImport(filePath)

	checkCallGraph(filePath, config.Rules.CallGraphRules)

	recoverBackupFile(filePath)

}

// parseCommandLineArgs parses the command line arguments and returns a RuntimeConfig object
func parseCommandLineArgs() *RuntimeConfig {
	config := RuntimeConfig{}
	flag.StringVar(&config.ProfileFilePath, "profile", "profiles/default.yaml", "")
	flag.StringVar(&config.SrcFilePath, "src", "", "")
	flag.Parse()
	fmt.Printf("Profile file: %s\n", config.ProfileFilePath)
	fmt.Printf("Source file: %s\n", config.SrcFilePath)

	return &config
}

func getSsaPkgs(codeFilePath string) (*ssa.Program, []*ssa.Package) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedImports | packages.NeedDeps,
	}

	// Load the packages

	pkgs, err := packages.Load(cfg, codeFilePath)
	if err != nil {
		fmt.Println("Error loading package:", err)
		return nil, nil
	}

	// Check for errors in the loaded packages
	if packages.PrintErrors(pkgs) > 0 {
		fmt.Println("Errors found in packages")
		return nil, nil
	}

	// // Filter out the specific package
	// var filteredPkgs []*packages.Package
	// for _, pkg := range pkgs {
	//     if pkg.PkgPath != "crypto/tls/fipsonly" {
	//         filteredPkgs = append(filteredPkgs, pkg)
	//     }
	// }

	// Build SSA packages
	prog, ssaPkgs := ssautil.AllPackages(pkgs, 0)
	if ssaPkgs == nil {
		fmt.Println("No SSA packages found")
		return nil, nil
	}

	prog.Build()
	return prog, ssaPkgs
}

func getSsaPkgs2(codeFilePath string) (*ssa.Program, []*ssa.Package) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedName,
		Tests: false,
	}

	// Load the packages
	pkgs, err := packages.Load(cfg, codeFilePath)
	if err != nil {
		fmt.Println("Error loading package:", err)
		return nil, nil
	}

	// Check for errors in the loaded packages
	if packages.PrintErrors(pkgs) > 0 {
		fmt.Println("Errors found in packages")
		return nil, nil
	}

	// Build SSA packages
	prog, ssaPkgs := ssautil.AllPackages(pkgs, 0)
	if ssaPkgs == nil {
		fmt.Println("No SSA packages found")
		return nil, nil
	}

	prog.Build()

	return prog, ssaPkgs
}

func checkCallGraph(codeFilePath string, rules []utils.Rule) {

	regexes := utils.GetCompiledRegexs(rules)

	// check if file exists
	if _, err := os.Stat(codeFilePath); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", codeFilePath)
		os.Exit(invalidBinaryExitCode)
	}

	// get directory of the code file
	dir := codeFilePath[:strings.LastIndex(codeFilePath, "/")]
	dir = path.Join(dir, "../../")
	// set the working directory to the directory of the code file
	err := os.Chdir(dir)

	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		os.Exit(3)
	}

	prog, ssaPkgs := getSsaPkgs(codeFilePath)

	_ = ssaPkgs

	// Generate the call graph

	cg := static.CallGraph(prog)
	for _, rule := range rules {
		match := false
		reg := regexes[rule.Name]
		for _, node := range cg.Nodes {
			for _, edge := range node.Out {
				// Get the called function
				calledFunction := edge.Callee.Func
				calledFuncName := calledFunction.String()

				if reg.MatchString(calledFuncName) {
					calledPos := prog.Fset.Position(calledFunction.Pos())
					fmt.Printf("Function: %s, File: %s, Line: %d\n", calledFuncName, calledPos.Filename, calledPos.Line)
					printCallTree(prog, cg, calledFunction, 0)
					match = true
					break
				}
			}
		}
		utils.PrintMatch(match, rule)
	}
}

func printCallTree(prog *ssa.Program, cg *callgraph.Graph, fn *ssa.Function, level int) {
    node := cg.Nodes[fn]
	if node == nil || level == 5{
		return
	}

    for _, edge := range node.In {
        caller := edge.Caller.Func
        fmt.Printf("%s%s\n", strings.Repeat("  ", level), caller.String())
		// print file and line number
		pos := prog.Fset.Position(caller.Pos())
		fmt.Printf("%s%s:%d\n", strings.Repeat("  ", level), pos.Filename, pos.Line)
        printCallTree(prog, cg, caller, level+1)
    }
}

func commentFipsImport(filePath string) error {
	// copy the file to a backup file
	backupFilePath := filePath + ".bak"

	// Open the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	// write the data to a backup file
	err = os.WriteFile(backupFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}
	new_data := strings.Replace(string(data), "_ \"crypto/tls/fipsonly\"", " // _ crypto/tls/fipsonly", -1)
	// write the new data to the file
	err = os.WriteFile(filePath, []byte(new_data), 0644)
	if err != nil {
		return fmt.Errorf("failed to write new data to file: %v", err)
	}
	return nil
}

// function to recover backup file and delete it
func recoverBackupFile(filePath string) error {
	backupFilePath := filePath + ".bak"
	data, err := os.ReadFile(filePath + ".bak")
	if err != nil {
		fmt.Printf("Error reading backup file: %v\n", err)
		os.Exit(5)
	}
	os.WriteFile(filePath, data, 0644)
	// delete the backup file
	err = os.Remove(backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to delete backup file: %v", err)
	}
	return nil
}

func findFunctionLineNumber(filePath, functionName string) (int, error) {
	// Create a new file set
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return 0, fmt.Errorf("could not parse file: %v", err)
	}

	// Traverse the AST to find the function declaration
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == functionName {
				// Get the line number of the function declaration
				position := fset.Position(funcDecl.Pos())
				return position.Line, nil
			}
		}
	}

	return 0, fmt.Errorf("function %s not found", functionName)
}
