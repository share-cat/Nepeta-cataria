package makejudge

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func RunCLI(args []string) {
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	genProblemDir := generateCmd.String("problem_dir", "", "Path to the problem directory (required). Example: judges/hello_world")
	genLang := generateCmd.String("lang", "", "Target language for the main solution template (python, cpp, javascript). If empty, uses the first language listed in problem's config.yaml.")

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

	default:
		log.Printf("Error: Unknown command '%s'\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: makejudge <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  generate    Generate boilerplate solution template and display config for a problem.")
	fmt.Println("              The generated files are intended for use by a judge server.")
	fmt.Println("\nOptions for generate:")
	fmt.Println("  -problem_dir string")
	fmt.Println("        Path to the problem directory (required). Example: judges/hello_world")
	fmt.Println("  -lang string")
	fmt.Println("        Target language for the main solution template (python, cpp, javascript).")
	fmt.Println("        If empty, uses the first language listed in the problem's config.yaml 'lang' field.")

}
