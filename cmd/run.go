package cmd

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/sankettank66/localmind/ollama"
	"github.com/sankettank66/localmind/ui"
)

const Version = "v1.1.0"

func printUsage() {
	// ui.PrintHeader()
	fmt.Println("\nAvailable flags:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  -%s \t %s (default: %v)\n", f.Name, f.Usage, f.DefValue)
	})
	fmt.Println()
}

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
	streamFlag := flag.Bool("stream", true, "Enable streaming responses")
	versionFlag := flag.Bool("version", false, "Show the application version")
	flag.Parse()

	// Show version and exit
	if *versionFlag {
		fmt.Printf("LocalMind %s\n", Version)
		return nil
	}

	// Show header on every start
	ui.PrintHeader()

	// Show available flags when no arguments provided
	if len(flag.Args()) > 0 || (flag.NFlag() == 0 && !*listFlag) {
		printUsage()
		return nil
	}

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

	// Print summary of what we're about to do
	fmt.Println()
	ui.Bold.Printf("  📝 Prompt  : %s\n", *promptFlag)
	ui.Bold.Printf("  🤖 Models  : %s\n", strings.Join(models, ", "))
	fmt.Println()

	if *streamFlag {
		streamChan := make(chan ollama.StreamEvent, 100)
		var wg sync.WaitGroup

		for _, model := range models {
			wg.Add(1)
			go func(m string) {
				defer wg.Done()
				ollama.Query(m, *promptFlag, streamChan)
			}(model)
		}

		go func() {
			wg.Wait()
			close(streamChan)
		}()

		ui.RenderStream(streamChan, models)
	} else {
		results, err := ollama.QueryAll(models, *promptFlag)
		if err != nil {
			return err
		}
		for _, r := range results {
			ui.PrintResult(r)
		}
	}

	return nil
}
