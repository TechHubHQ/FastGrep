# ⚡ FastGrep

A blazing-fast, concurrent text search tool written in Go. FastGrep leverages a multi-level worker pool architecture to search for patterns across files and directories with high throughput — designed as a performant alternative to traditional `grep` for large-scale searches.

## Features

- **Multi-level Concurrency** — Two-tier worker pool (file-level + chunk-level) for maximum parallelism
- **Large File Support** — Files are split into 1 MB chunks and searched in parallel across CPU cores
- **Buffer Pooling** — `sync.Pool`-based buffer reuse to minimize memory allocations and GC pressure
- **Glob Pattern Matching** — Supports wildcards (`*`) and multi-file argument expansion
- **Recursive Directory Search** — Traverse entire directory trees with the `-r` flag
- **Case-Insensitive Search** — Optional case-insensitive matching with `-i`
- **Colorized Output** — Highlighted filenames (blue), line numbers (red), and matched patterns (yellow/bold)
- **Cross-Platform** — Works on both Windows and Linux

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) 1.26+ installed

### Build from Source

```bash
git clone https://github.com/TechHubHQ/FastGrep.git
cd FastGrep
go build -o fastgrep ./cmd/fastgrep
```

The compiled binary will be created in the current directory.

## Usage

```bash
fastgrep [OPTIONS] <search-pattern> <file-or-directory>...
```

### Options

| Flag | Description                      |
|------|----------------------------------|
| `-i` | Ignore case in search            |
| `-r` | Search directories recursively   |
| `-h` | Show help message                |

### Examples

**Search in a single file:**

```bash
fastgrep "hello" main.go
```

**Case-insensitive search across log files:**

```bash
fastgrep -i "error" ./logs/*
```

**Recursive search in all files under current directory:**

```bash
fastgrep -r "TODO" .
```

**Search across multiple specific files:**

```bash
fastgrep "func main" server.go handler.go router.go
```

### Output Format

Results are printed in a colorized `grep`-style format:

```text
<filename>:<line-number>:<matching-line>
```

- **Filename** — displayed in blue
- **Line number** — displayed in red
- **Matched pattern** — highlighted in bold yellow within the line

## Architecture

FastGrep uses a **producer-consumer model** with two levels of worker pools to maximize throughput:

```text
┌─────────────┐
│ Orchestrator│
└──────┬──────┘
       │  dispatches files
       ▼
┌──────────────────┐
│ File Worker Pool │  N workers (min of file count, CPU cores)
│  ┌────┐ ┌────┐   │
│  │ FW │ │ FW │…  │  Each reads a file in 1 MB chunks
│  └──┬─┘ └──┬─┘   │
└─────┼───────┼────┘
      │       │  dispatches search jobs
      ▼       ▼
┌──────────────────┐
│ Chunk Worker Pool│  N workers (CPU core count)
│  ┌────┐ ┌────┐   │
│  │ CW │ │ CW │   │  Each searches a chunk for the pattern
│  └──┬─┘ └──┬─┘   │
└─────┼──────┼─────┘
      │      │  matched results
      ▼      ▼
┌──────────────────┐
│  Printer (stdout)│  Single goroutine — serialized, no interleaving
└──────────────────┘
```

### Key Design Decisions

| Concern                | Approach                                                                                                         |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------- |
| **Chunk boundaries**   | After reading a 1 MB chunk, the reader advances to the next `\n` to ensure lines are never split across chunks   |
| **Memory efficiency**  | A `sync.Pool` of reusable 1 MB byte buffers reduces allocations; chunks are copied to right-sized slices         |
| **Output correctness** | All results flow through a shared channel to a dedicated printer goroutine, eliminating interleaved output       |
| **Binary file safety** | `.exe` files are automatically skipped during file discovery                                                     |

## Project Structure

```text
FastGrep/
├── cmd/
│   └── fastgrep/
│       └── main.go              # Entry point
├── internal/
│   ├── cli/
│   │   ├── parser.go            # CLI argument parsing & help text
│   │   └── run.go               # Top-level run function
│   ├── config/
│   │   └── config.go            # Constants (chunk size, etc.)
│   ├── fsutil/
│   │   └── file_utils.go        # Glob expansion, file filtering, directory walking
│   ├── orchestrator/
│   │   └── orchestrator.go      # Wires up worker pools, channels, and printer
│   ├── search/
│   │   └── matcher.go           # Pattern matching logic, Match/SearchJob types
│   └── worker/
│       ├── buffer_pool.go       # sync.Pool-based buffer management
│       ├── chunk_worker.go      # Chunk-level worker pool
│       └── file_worker.go       # File-level worker pool with chunked reading
├── benchmark_test.go            # End-to-end benchmarks (single & multi-file)
├── docs/
│   ├── cli-plan.md              # CLI design notes
│   └── logic-plan.md            # Concurrency design notes
├── go.mod
└── .gitignore
```

## Benchmarks

FastGrep includes end-to-end benchmarks for single-file and multi-file scenarios:

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Single 50 MB file benchmark
go test -bench=BenchmarkEndToEnd -benchmem

# Multi-file benchmark (5 × 10 MB)
go test -bench=BenchmarkEndToEndMultipleFiles -benchmem
```

## Roadmap

- [ ] Regex pattern support
- [ ] Context lines (`-A`, `-B`, `-C` flags)
- [ ] File type filtering (`--include`, `--exclude`)
- [ ] Match count summary (`-c` flag)
- [ ] JSON output format

## 📄 License

This project is open source.
