
```text
  ██╗      ██████╗  ██████╗ █████╗ ██╗      ███╗   ███╗██╗███╗   ██╗██████╗
  ██║     ██╔═══██╗██╔════╝██╔══██╗██║      ████╗ ████║██║████╗  ██║██╔══██╗
  ██║     ██║   ██║██║     ███████║██║      ██╔████╔██║██║██╔██╗ ██║██║  ██║
  ██║     ██║   ██║██║     ██╔══██║██║      ██║╚██╔╝██║██║██║╚██╗██║██║  ██║
  ███████╗╚██████╔╝╚██████╗██║  ██║███████╗ ██║ ╚═╝ ██║██║██║ ╚████║██████╔╝
  ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═════╝
  ```

> Run the same prompt across multiple local LLMs simultaneously - compare outputs, speed, and latency. All on your machine. No API keys. No cloud.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
![Ollama](https://img.shields.io/badge/Ollama-required-black?style=flat)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=flat)

---

## The Problem

You have 4 local models installed. You don't know which one is best for your task. So you open Ollama, run one model, copy the output, switch models, run again... and repeat.

**LocalMind fixes that.** One command. All models. Side by side.

---

## Demo

```bash
$ localmind -models "llama3,mistral,phi3" -prompt "explain recursion in one line"

  MODEL: llama3
  ──────────────────────────────────────────────────────────
  A function that calls itself until a base condition is met.
  ⏱  TTFT: 340ms  |  Total: 1.2s

  MODEL: mistral
  ──────────────────────────────────────────────────────────
  Recursion is when a function solves a problem by calling itself on smaller inputs.
  ⏱  TTFT: 280ms  |  Total: 980ms

  MODEL: phi3
  ──────────────────────────────────────────────────────────
  A function that repeatedly calls itself with a simpler version of the original problem.
  ⏱  TTFT: 190ms  |  Total: 720ms
```

---

## Features

- **Parallel execution** - all models run simultaneously, not one after another
- **Real metrics** - Time to First Token (TTFT) and total response time per model
- **Model discovery** - list all your locally installed Ollama models
- **Zero dependencies** - no API keys, no internet, runs fully local
- **Clean output** - color-coded terminal UI that's easy to read

---

## Requirements

- [Go 1.25+](https://golang.org/dl/)
- [Ollama](https://ollama.com/) running locally with at least one model pulled

---

## Installation

### Option 1 - Build from source

```bash
git clone https://github.com/sankettank66/localmind.git
cd localmind
go build -o localmind .
```

**Windows:**
```bash
go build -o localmind.exe .
```

### Option 2 - Install with Go

```bash
go install github.com/sankettank66/localmind@latest
```

---

## Usage

### List available models

```bash
localmind -list
```

### Compare models on a prompt

```bash
localmind -models "llama3,mistral" -prompt "what is a goroutine?"
```

### Compare more models

```bash
localmind -models "llama3,mistral,phi3,gemma" -prompt "write a binary search in python"
```

---

## How It Works

1. LocalMind connects to your local Ollama instance at `localhost:11434`
2. It sends your prompt to all specified models **in parallel** using goroutines
3. Results stream back as each model finishes
4. TTFT and total time are captured per model so you can compare speed

---

## Project Structure

```
localmind/
├── main.go           # Entry point
├── cmd/
│   └── run.go        # CLI flags and orchestration
├── ollama/
│   └── client.go     # Ollama API client
└── ui/
    └── display.go    # Terminal output formatting
```

---

## Roadmap

- [x] Parallel model querying
- [x] TTFT and total time metrics
- [x] Model listing
- [ ] Real-time streaming output (token by token)
- [ ] Web dashboard (React)
- [ ] Custom benchmark suites (coding, reasoning, math)
- [ ] Export results as markdown/JSON report
- [ ] Compare local models vs cloud APIs
- [ ] ELO scoring based on your personal history

---

## Contributing

Contributions are welcome! Please open an issue first to discuss what you'd like to change.

```bash
git clone https://github.com/sankettank66/localmind.git
cd localmind
go mod tidy
go run . -list
```

---

> Built because switching between Ollama models manually is painful. LocalMind makes it one command.