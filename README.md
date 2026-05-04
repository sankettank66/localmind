
```text
  ██╗      ██████╗  ██████╗ █████╗ ██╗      ███╗   ███╗██╗███╗   ██╗██████╗
  ██║     ██╔═══██╗██╔════╝██╔══██╗██║      ████╗ ████║██║████╗  ██║██╔══██╗
  ██║     ██║   ██║██║     ███████║██║      ██╔████╔██║██║██╔██╗ ██║██║  ██║
  ██║     ██║   ██║██║     ██╔══██║██║      ██║╚██╔╝██║██║██║╚██╗██║██║  ██║
  ███████╗╚██████╔╝╚██████╗██║  ██║███████╗ ██║ ╚═╝ ██║██║██║ ╚████║██████╔╝
  ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═════╝
  ```

> Run the same prompt across multiple local LLMs simultaneously - compare outputs, speed, and latency. All on your machine. No API keys. No cloud.

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go&logoColor=white)
![Ollama](https://img.shields.io/badge/Ollama-required-black?style=flat)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=flat)

---

## The Problem

You have multiple local models installed and you're not sure which one provides the best reasoning or fastest response for a specific task. Manually switching and re-running prompts is tedious.

**LocalMind fixes that.** One command. All models. Side-by-side real-time comparison.

---

## Features

- **🚀 Real-time Streaming** - Watch tokens appear live from all models simultaneously.
- **📦 Sophisticated Boxed UI** - Clean, modern terminal layout with automatic word wrapping.
- **🔄 Service Lifecycle Management** - Automatically checks if Ollama is installed and starts the service if it's not running.
- **📊 Deep Metrics** - Detailed stats per model:
    - **IN/OUT:** Token counts for prompt and generation.
    - **TTFT:** Time to First Token.
    - **LOAD:** Model load duration.
    - **TOTAL:** Complete execution time.
- **Zero Dependencies** - Fully local, no cloud, no API keys.

---

## Installation

### Build from source

```bash
git clone https://github.com/sankettank66/localmind.git
cd localmind
go build -o localmind .
```

---

## Usage

### Compare models with real-time streaming (Default)

```bash
localmind -models "llama3,mistral" -prompt "explain quantum entanglement"
```

### Disable streaming for batch output

```bash
localmind -models "llama3,mistral" -prompt "write a go function to sort a slice" -stream=false
```

### List available models

```bash
localmind -list
```

---

## How It Works

1. **Pre-flight Checks**: LocalMind verifies your Ollama installation and starts the background service if necessary.
2. **Parallel Orchestration**: Uses Go routines to query multiple models concurrently.
3. **Reactive UI**: Implements a custom ANSI-based rendering engine that updates multiple "boxes" on your terminal as streams arrive.
4. **Metadata Capture**: Aggregates precise timing and token data from Ollama's response headers.

---

## Roadmap

- [x] Parallel model querying
- [x] TTFT and total time metrics
- [x] Model listing
- [x] Real-time streaming output
- [x] Service auto-start & lifecycle management
- [ ] Web dashboard (React)
- [ ] Custom benchmark suites (coding, reasoning, math)
- [ ] Export results as markdown/JSON report
- [ ] ELO scoring based on your personal history

---

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

```bash
go run . -list
```

---

> Built because switching between Ollama models manually is painful. LocalMind makes it one command.
