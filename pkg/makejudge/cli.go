package makejudge

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"share-cat/Nepeta-cataria/pkg/loadjudge"
)

func printUsage() {
	fmt.Println("makejudge CLI Usage:")
	fmt.Println("  makejudge generate -problem_dir <path> [-lang <language>]")
	fmt.Println("  makejudge testgen -problem_dir <path> -code_file <path> -lang <language> [-output <path>]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  generate: Generate solution template files for a problem")
	fmt.Println("  testgen:  Generate executable code with embedded test data")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  makejudge generate -problem_dir judges/hello_world -lang python")
	fmt.Println("  makejudge testgen -problem_dir judges/hello_world -code_file solution.py -lang python")
	fmt.Println("  makejudge testgen -problem_dir judges/hello_world -code_file solution.py -lang python -output test_runner.py")
}

func RunCLI(args []string) {
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	genProblemDir := generateCmd.String("problem_dir", "", "Path to the problem directory (required). Example: judges/hello_world")
	genLang := generateCmd.String("lang", "", "Target language for the main solution template (python, cpp, javascript). If empty, uses the first language listed in problem's config.yaml.")

	testgenCmd := flag.NewFlagSet("testgen", flag.ExitOnError)
	testProblemDir := testgenCmd.String("problem_dir", "", "Path to the problem directory (required)")
	testCodeFile := testgenCmd.String("code_file", "", "Path to the user code file (required)")
	testLang := testgenCmd.String("lang", "", "Programming language (python, cpp, javascript) (required)")
	testOutput := testgenCmd.String("output", "", "Output file path (optional). If not specified, prints to stdout")

	log.SetFlags(log.LstdFlags)

	switch args[0] {
	case "generate":
		err := generateCmd.Parse(args[1:])
		if err != nil {
			log.Fatalf("Error parsing 'generate' arguments: %v", err)
		}

		if *genProblemDir == "" {
			log.Println("Error: -problem_dir is required for generate command.")
			generateCmd.Usage()
			os.Exit(1)
		}
		absProblemDir, err := filepath.Abs(*genProblemDir)
		if err != nil {
			log.Fatalf("Error getting absolute path for problem_dir: %v", err)
		}

		if _, errStat := os.Stat(absProblemDir); os.IsNotExist(errStat) {
			log.Fatalf("Error: Problem directory '%s' does not exist.", absProblemDir)
		}

		log.Printf("Executing 'generate' for problem directory: %s, target template language: %s", absProblemDir, *genLang)
		err = ScaffoldProblem(absProblemDir, *genLang)
		if err != nil {
			log.Fatalf("Error during 'generate' command: %v", err)
		}
		log.Println("'generate' command completed successfully.")

	case "testgen":
		err := testgenCmd.Parse(args[1:])
		if err != nil {
			log.Fatalf("Error parsing 'testgen' arguments: %v", err)
		}

		if *testProblemDir == "" {
			log.Println("Error: -problem_dir is required for testgen command.")
			testgenCmd.Usage()
			os.Exit(1)
		}

		if *testCodeFile == "" {
			log.Println("Error: -code_file is required for testgen command.")
			testgenCmd.Usage()
			os.Exit(1)
		}

		if *testLang == "" {
			log.Println("Error: -lang is required for testgen command.")
			testgenCmd.Usage()
			os.Exit(1)
		}

		err = runTestGen(*testProblemDir, *testCodeFile, *testLang, *testOutput)
		if err != nil {
			log.Fatalf("Error during 'testgen' command: %v", err)
		}

	default:
		log.Printf("Error: Unknown command '%s'\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

func runTestGen(problemDir, codeFile, language, outputFile string) error {
	absProblemDir, err := filepath.Abs(problemDir)
	if err != nil {
		return fmt.Errorf("getting absolute path for problem_dir: %w", err)
	}

	absCodeFile, err := filepath.Abs(codeFile)
	if err != nil {
		return fmt.Errorf("getting absolute path for code_file: %w", err)
	}

	if _, err := os.Stat(absProblemDir); os.IsNotExist(err) {
		return fmt.Errorf("problem directory '%s' does not exist", absProblemDir)
	}

	if _, err := os.Stat(absCodeFile); os.IsNotExist(err) {
		return fmt.Errorf("code file '%s' does not exist", absCodeFile)
	}

	log.Printf("Generating executable code for problem: %s", absProblemDir)
	log.Printf("Using code file: %s", absCodeFile)
	log.Printf("Target language: %s", language)

	testDataconfig, err := loadjudge.LoadTestDataConfig(absProblemDir)
	if err != nil {
		return fmt.Errorf("loading problem config: %w", err)
	}

	userCodeBytes, err := os.ReadFile(absCodeFile)
	if err != nil {
		return fmt.Errorf("reading user code file: %w", err)
	}
	userCode := string(userCodeBytes)

	for _, testcase := range testDataconfig.TestCases {

		executableCode, err := GenerateExecutableCodeWithTestData(testDataconfig, testcase, userCode, language)
		if err != nil {
			return fmt.Errorf("generating executable code: %w", err)
		}
		if outputFile != "" {
			absOutputFile, err := filepath.Abs(testcase.Name + "_" + outputFile)
			if err != nil {
				return fmt.Errorf("getting absolute path for output file: %w", err)
			}

			err = os.WriteFile(absOutputFile, []byte(executableCode), 0644)
			if err != nil {
				return fmt.Errorf("writing output file: %w", err)
			}

			log.Printf("Executable code written to: %s", absOutputFile)
		} else {
			fmt.Println("=== Generated Executable Code ===")
			fmt.Println(executableCode)
		}
	}

	log.Println("'testgen' command completed successfully.")
	return nil
}
