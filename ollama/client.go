package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

type ModelResult struct {
	Model              string
	Response           string
	PromptEvalCount    int
	PromptEvalDuration time.Duration
	EvalCount          int
	EvalDuration       time.Duration
	LoadDuration       time.Duration
	TotalTime          time.Duration
	Error              error
}

type StreamEventType int

const (
	EventToken StreamEventType = iota
	EventDone
	EventError
)

type StreamEvent struct {
	Type               StreamEventType
	Model              string
	Token              string
	PromptEvalCount    int
	PromptEvalDuration time.Duration
	EvalCount          int
	EvalDuration       time.Duration
	LoadDuration       time.Duration
	TotalTime          time.Duration
	Err                error
}

func IsOllamaInstalled() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

func StartOllama() error {
	cmd := exec.Command("ollama", "serve")

	// Detach from current process
	err := cmd.Start()
	if err != nil {
		return err
	}

	fmt.Println("✅ Ollama started in background")
	return nil
}

func WaitForOllama() bool {
	for i := 0; i < 5; i++ {
		if IsOllamaRunning() {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func IsOllamaRunning() bool {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func ListModels() ([]string, error) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, fmt.Errorf("cannot reach Ollama - is it running? %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	names := make([]string, len(result.Models))
	for i, m := range result.Models {
		names[i] = m.Name
	}
	return names, nil
}

func Query(model, prompt string, resultChan chan<- StreamEvent) {
	reqBody, err := json.Marshal(GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	})
	if err != nil {
		resultChan <- StreamEvent{Type: EventError, Model: model, Err: err}
		return
	}

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		resultChan <- StreamEvent{Type: EventError, Model: model, Err: err}
		return
	}
	defer resp.Body.Close()

	ttftSet := false
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var chunk GenerateResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			continue
		}

		if chunk.Done {
			resultChan <- StreamEvent{
				Type:               EventDone,
				Model:              model,
				PromptEvalCount:    chunk.PromptEvalCount,
				PromptEvalDuration: time.Duration(chunk.PromptEvalDuration),
				EvalCount:          chunk.EvalCount,
				EvalDuration:       time.Duration(chunk.EvalDuration),
				LoadDuration:       time.Duration(chunk.LoadDuration),
				TotalTime:          time.Duration(chunk.TotalDuration),
			}
			return
		}

		if chunk.Response == "" {
			continue
		}

		if !ttftSet {
			ttftSet = true
		}
		resultChan <- StreamEvent{Type: EventToken, Model: model, Token: chunk.Response}
	}

	if err := scanner.Err(); err != nil {
		resultChan <- StreamEvent{Type: EventError, Model: model, Err: err}
		return
	}
}

func QueryAll(models []string, prompt string) ([]ModelResult, error) {
	results := make([]ModelResult, len(models))
	var wg sync.WaitGroup

	for i, model := range models {
		wg.Add(1)
		go func(idx int, m string) {
			defer wg.Done()
			result := QuerySync(m, prompt)
			results[idx] = result
		}(i, model)
	}

	wg.Wait()
	return results, nil
}

func QuerySync(model, prompt string) ModelResult {
	reqBody, err := json.Marshal(GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return ModelResult{Model: model, Error: err}
	}

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return ModelResult{Model: model, Error: err}
	}
	defer resp.Body.Close()

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ModelResult{Model: model, Error: err}
	}

	return ModelResult{
		Model:              model,
		Response:           result.Response,
		PromptEvalCount:    result.PromptEvalCount,
		PromptEvalDuration: time.Duration(result.PromptEvalDuration),
		EvalCount:          result.EvalCount,
		EvalDuration:       time.Duration(result.EvalDuration),
		LoadDuration:       time.Duration(result.LoadDuration),
		TotalTime:          time.Duration(result.TotalDuration),
	}
}
