package cmd

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/sankettank66/localmind/ollama"
	"github.com/sankettank66/localmind/ui"
)

func ensureOllama() error {
	ui.ShowLoading("Checking Ollama installation...")
	if !ollama.IsOllamaInstalled() {
		ui.UpdateLoading(false, "Ollama is not installed. Please install it from ollama.com")
		return fmt.Errorf("ollama not found")
	}
	ui.UpdateLoading(true, "Ollama is installed")

	ui.ShowLoading("Checking if Ollama is running...")
	if !ollama.IsOllamaRunning() {
		ui.UpdateLoading(false, "Ollama is not running. Starting it now...")

		ui.ShowLoading("Starting Ollama service...")
		err := ollama.StartOllama()
		if err != nil {
			ui.UpdateLoading(false, "Failed to start Ollama")
			return err
		}

		if !ollama.WaitForOllama() {
			ui.UpdateLoading(false, "Ollama service timed out")
			return fmt.Errorf("ollama timeout")
		}
		ui.UpdateLoading(true, "Ollama service started")
	} else {
		ui.UpdateLoading(true, "Ollama service is active")
	}
	return nil
}

func Execute() error {
	modelsFlag := flag.String("models", "", "Comma-separated list of models (e.g. llama3,mistral)")
	promptFlag := flag.String("prompt", "", "Prompt to send to all models")
	listFlag := flag.Bool("list", false, "List all available Ollama models")
	flag.Parse()

	// Ensure Ollama is ready before doing anything else
	if err := ensureOllama(); err != nil {
		return err
	}

	// List models
	if *listFlag {
		models, err := ollama.ListModels()
		if err != nil {
			return err
		}
		fmt.Println("\nAvailable models:")
		for _, m := range models {
			fmt.Printf("  • %s\n", m)
		}
		return nil
	}

	// Validate flags
	if *promptFlag == "" {
		return fmt.Errorf("please provide a prompt using -prompt \"your question\"")
	}
	if *modelsFlag == "" {
		return fmt.Errorf("please provide models using -models \"llama3,mistral\"")
	}

	models := strings.Split(*modelsFlag, ",")
	for i := range models {
		models[i] = strings.TrimSpace(models[i])
	}

	ui.PrintHeader(*promptFlag, models)

	// Run all models in parallel
	resultChan := make(chan ollama.ModelResult, len(models))
	var wg sync.WaitGroup

	for _, model := range models {
		wg.Add(1)
		go func(m string) {
			defer wg.Done()
			ollama.Query(m, *promptFlag, resultChan)
		}(model)
	}

	// Close channel when all done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Print results as they come in
	for result := range resultChan {
		ui.PrintResult(result)
	}

	return nil
}
