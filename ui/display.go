package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/sankettank66/localmind/ollama"
)

var (
	Success = color.New(color.FgCyan, color.Bold)
	Bold    = color.New(color.Bold)
	Dim     = color.New(color.FgHiBlack)
	Metric  = color.New(color.FgMagenta)
	Info    = color.New(color.FgBlue)
	Label   = color.New(color.FgHiYellow)
)

type modelResult struct {
	model              string
	promptEvalCount    int
	promptEvalDuration time.Duration
	evalCount          int
	evalDuration       time.Duration
	loadDuration       time.Duration
	totalTime          time.Duration
	response           strings.Builder
	err                error
}

func RenderStream(streamChan <-chan ollama.StreamEvent, models []string) {
	numModels := len(models)
	states := make(map[string]*modelResult)
	for _, m := range models {
		states[m] = &modelResult{model: m}
	}

	var mu sync.Mutex
	doneCount := 0

	// Initial print
	mu.Lock()
	for _, m := range models {
		printBox(states[m])
	}
	mu.Unlock()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	updatesNeeded := false
	currentTotalLines := 0
	for _, m := range models {
		currentTotalLines += countBoxLines(states[m])
	}

	for {
		select {
		case event, ok := <-streamChan:
			if !ok {
				goto finalRedraw
			}
			mu.Lock()
			s := states[event.Model]
			switch event.Type {
			case ollama.EventToken:
				s.response.WriteString(event.Token)
			case ollama.EventDone:
				s.promptEvalCount = event.PromptEvalCount
				s.promptEvalDuration = event.PromptEvalDuration
				s.evalCount = event.EvalCount
				s.evalDuration = event.EvalDuration
				s.loadDuration = event.LoadDuration
				s.totalTime = event.TotalTime
				doneCount++
			case ollama.EventError:
				s.err = event.Err
				doneCount++
			}
			updatesNeeded = true
			mu.Unlock()

		case <-ticker.C:
			if updatesNeeded {
				mu.Lock()
				// 1. Move up by the lines we PREVIOUSLY printed
				moveUp(currentTotalLines)

				// 2. Redraw all boxes
				newTotalLines := 0
				for _, m := range models {
					printBox(states[m])
					newTotalLines += countBoxLines(states[m])
				}

				// 3. Update the line count for next moveUp
				currentTotalLines = newTotalLines
				updatesNeeded = false
				mu.Unlock()
			}
			if doneCount == numModels {
				goto finalRedraw
			}
		}
	}

finalRedraw:
	mu.Lock()
	moveUp(currentTotalLines)
	for _, m := range models {
		printBox(states[m])
	}
	mu.Unlock()
}

func countBoxLines(s *modelResult) int {
	width := 70
	fixedOverhead := 7

	content := strings.TrimSpace(s.response.String())
	if content == "" {
		content = "Thinking..."
	}

	var contentLines []string
	if s.err != nil {
		contentLines = wrapText(fmt.Sprintf("Error: %s", s.err), width-4)
	} else {
		contentLines = wrapText(content, width-4)
	}

	return fixedOverhead + len(contentLines)
}

func moveUp(lines int) {
	if lines > 0 {
		fmt.Printf("\033[%dA", lines)
	}
}

func printBox(s *modelResult) {
	width := 70
	line := strings.Repeat("─", width)

	// 1: Top border
	fmt.Print("\033[K")
	fmt.Printf("%s\n", Dim.Sprint("╭"+line+"╮"))

	// 2: Title & Status
	status := Success.Sprint("RUNNING...")
	plainStatus := "RUNNING..."
	if s.err != nil {
		status = color.RedString("FAILED")
		plainStatus = "FAILED"
	} else if s.totalTime > 0 {
		status = Success.Sprint("COMPLETE")
		plainStatus = "COMPLETE"
	}

	titleText := fmt.Sprintf(" %s %s", Info.Sprint("●"), Bold.Sprint(s.model))
	plainTitle := fmt.Sprintf(" ● %s", s.model)

	paddingCount := width - len(plainTitle) - len(plainStatus) - 2
	if paddingCount < 1 {
		paddingCount = 1
	}

	fmt.Print("\033[K")
	fmt.Printf("%s%s%s%s%s\n",
		Dim.Sprint("│"),
		titleText,
		strings.Repeat(" ", paddingCount),
		status,
		Dim.Sprint(" │"),
	)

	// 3: Separator 1
	fmt.Print("\033[K")
	fmt.Printf("%s\n", Dim.Sprint("├"+line+"┤"))

	// 4: Content
	var contentLines []string
	if s.err != nil {
		contentLines = wrapText(fmt.Sprintf("Error: %s", s.err), width-4)
	} else {
		content := strings.TrimSpace(s.response.String())
		if content == "" {
			content = "Thinking..."
		}
		contentLines = wrapText(content, width-4)
	}

	for _, l := range contentLines {
		fmt.Print("\033[K")
		fmt.Printf("%s  %-66s  %s\n", Dim.Sprint("│"), l, Dim.Sprint("│"))
	}

	// 5: Separator 2
	fmt.Print("\033[K")
	fmt.Printf("%s\n", Dim.Sprint("├"+line+"┤"))

	// 6: Stats Section
	fmt.Print("\033[K")
	if s.totalTime > 0 {
		statsLine := fmt.Sprintf(" %s%d %s%d %s%s %s%s %s%s",
			Label.Sprint("IN:"), s.promptEvalCount,
			Label.Sprint("OUT:"), s.evalCount,
			Label.Sprint("TTFT:"), Metric.Sprint(s.promptEvalDuration.Round(time.Millisecond)),
			Label.Sprint("LOAD:"), Metric.Sprint(s.loadDuration.Round(time.Millisecond)),
			Label.Sprint("TOTAL:"), Metric.Sprint(s.totalTime.Round(time.Millisecond)),
		)
		plainStats := fmt.Sprintf(" IN:%d OUT:%d TTFT:%s LOAD:%s TOTAL:%s",
			s.promptEvalCount, s.evalCount,
			s.promptEvalDuration.Round(time.Millisecond),
			s.loadDuration.Round(time.Millisecond),
			s.totalTime.Round(time.Millisecond),
		)
		statsPadding := width - len(plainStats) - 1
		if statsPadding < 0 {
			statsPadding = 0
		}
		fmt.Printf("%s%s%s%s\n", Dim.Sprint("│"), statsLine, strings.Repeat(" ", statsPadding), Dim.Sprint("│"))
	} else {
		fmt.Printf("%s %-68s %s\n", Dim.Sprint("│"), Info.Sprint(" Streaming response..."), Dim.Sprint("│"))
	}

	// 7: Bottom border
	fmt.Print("\033[K")
	fmt.Printf("%s\n", Dim.Sprint("╰"+line+"╯"))

	// 8: Trailing spacer newline
	fmt.Print("\033[K")
	fmt.Println()
}

func PrintResult(r ollama.ModelResult) {
	s := &modelResult{
		model:              r.Model,
		promptEvalCount:    r.PromptEvalCount,
		promptEvalDuration: r.PromptEvalDuration,
		evalCount:          r.EvalCount,
		evalDuration:       r.EvalDuration,
		loadDuration:       r.LoadDuration,
		totalTime:          r.TotalTime,
		err:                r.Error,
	}
	s.response.WriteString(r.Response)
	printBox(s)
}

func wrapText(text string, width int) []string {
	var lines []string
	paragraphs := strings.Split(text, "\n")

	for _, p := range paragraphs {
		words := strings.Fields(p)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		currentLine := ""
		for _, word := range words {
			if len(currentLine)+len(word)+1 > width {
				if currentLine != "" {
					lines = append(lines, currentLine)
				}
				currentLine = word
			} else {
				if currentLine != "" {
					currentLine += " "
				}
				currentLine += word
			}
		}
		if currentLine != "" {
			lines = append(lines, currentLine)
		}
	}
	return lines
}

func PrintHeader() {
	fmt.Println()
	Success.Println("  ██╗      ██████╗  ██████╗ █████╗ ██╗      ███╗   ███╗██╗███╗   ██╗██████╗")
	Success.Println("  ██║     ██╔═══██╗██╔════╝██╔══██╗██║      ████╗ ████║██║████╗  ██║██╔══██╗")
	Success.Println("  ██║     ██║   ██║██║     ███████║██║      ██╔████╔██║██║██╔██╗ ██║██║  ██║")
	Success.Println("  ██║     ██║   ██║██║     ██╔══██║██║      ██║╚██╔╝██║██║██║╚██╗██║██║  ██║")
	Success.Println("  ███████╗╚██████╔╝╚██████╗██║  ██║███████╗ ██║ ╚═╝ ██║██║██║ ╚████║██████╔╝")
	Success.Println("  ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═════╝")
	fmt.Println()
}

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

func PrintError(err error) {
	fmt.Printf("\n  %s %s\n\n", color.RedString("✘ Error:"), err.Error())
}
