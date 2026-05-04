package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

type ModelResult struct {
	Model     string
	Response  string
	TTFT      time.Duration // time to first token
	TotalTime time.Duration
	Error     error
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

func Query(model, prompt string, resultChan chan<- ModelResult) {
	start := time.Now()

	reqBody, _ := json.Marshal(GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		resultChan <- ModelResult{Model: model, Error: err}
		return
	}
	defer resp.Body.Close()

	ttft := time.Since(start)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultChan <- ModelResult{Model: model, Error: err}
		return
	}

	var genResp GenerateResponse
	if err := json.Unmarshal(body, &genResp); err != nil {
		resultChan <- ModelResult{Model: model, Error: err}
		return
	}

	resultChan <- ModelResult{
		Model:     model,
		Response:  genResp.Response,
		TTFT:      ttft,
		TotalTime: time.Since(start),
	}
}
