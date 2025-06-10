package main

import (
	"fmt"
	"log"
	"share-cat/Nepeta-cataria/pkg/catwheel"
)

func main() {
	client := catwheel.NewClient("http://localhost:3200")

	options := catwheel.DefaultExecuteOptions().
		WithLanguage(catwheel.LanguagePython).
		WithTimeout(10).
		WithMemoryLimit(128 * 1024 * 1024)

	log.Printf("Executing code with options: %+v\n", options)

	response, err := client.ExecuteWithOptions("print('hello world')", options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Execution ID: %s\n", response.ExecutionID)
	fmt.Printf("Status: %s\n", response.Result.Status)
	fmt.Printf("Output: %s\n", response.Result.Stdout)
	fmt.Printf("Error: %s\n", response.Result.Stderr)

	execution, err := client.GetExecution(response.ExecutionID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Execution time: %v\n", execution.Result.Duration())

	executions, err := client.ListExecutions()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total executions: %d\n", len(executions.Executions))

	cleanup, err := client.CleanupExecution(response.ExecutionID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cleanup message: %s\n", cleanup.Message)

	health, err := client.HealthCheck()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Health status: %s\n", health.Status)

}
