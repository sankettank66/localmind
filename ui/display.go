package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/sankettank66/localmind/ollama"
)

var (
	header   = color.New(color.FgCyan, color.Bold)
	success  = color.New(color.FgGreen, color.Bold)
	metric   = color.New(color.FgYellow)
	errColor = color.New(color.FgRed)
	Info     = color.New(color.FgBlue)
)

func ShowLoading(message string) {
	fmt.Printf("  %s %-30s", Info.Sprint("⌛"), message)
}

func UpdateLoading(success bool, message string) {
	fmt.Print("\r")
	if success {
		fmt.Printf("  %s %-30s\n", color.GreenString("✅"), message)
	} else {
		fmt.Printf("  %s %-30s\n", color.RedString("❌"), message)
	}
}

func PrintResult(r ollama.ModelResult) {
	fmt.Println(strings.Repeat("─", 60))
	header.Printf("  MODEL: %s\n", r.Model)
	fmt.Println(strings.Repeat("─", 60))

	if r.Error != nil {
		errColor.Printf("  ERROR: %s\n\n", r.Error)
		return
	}

	fmt.Printf("  %s\n\n", r.Response)
	metric.Printf("  ⏱  TTFT: %s\n", r.TTFT.Round(1000000))
	metric.Printf("  ⏱  Total: %s\n\n", r.TotalTime.Round(1000000))
}

func PrintHeader(prompt string, models []string) {
	fmt.Println()
	success.Println("  ██╗      ██████╗  ██████╗ █████╗ ██╗      ███╗   ███╗██╗███╗   ██╗██████╗")
	success.Println("  ██║     ██╔═══██╗██╔════╝██╔══██╗██║      ████╗ ████║██║████╗  ██║██╔══██╗")
	success.Println("  ██║     ██║   ██║██║     ███████║██║      ██╔████╔██║██║██╔██╗ ██║██║  ██║")
	success.Println("  ██║     ██║   ██║██║     ██╔══██║██║      ██║╚██╔╝██║██║██║╚██╗██║██║  ██║")
	success.Println("  ███████╗╚██████╔╝╚██████╗██║  ██║███████╗ ██║ ╚═╝ ██║██║██║ ╚████║██████╔╝")
	success.Println("  ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═════╝")
	fmt.Println()
	header.Printf("  Prompt  : %s\n", prompt)
	header.Printf("  Models  : %s\n", strings.Join(models, ", "))
	fmt.Println()
}
