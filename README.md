# BrowserForge

BrowserForge is a Go library for generating realistic browser fingerprints and HTTP headers for web scraping and testing applications.

## Features

- Generate consistent browser fingerprints with matching HTTP headers
- Based on Bayesian networks for realistic, correlated values
- Lightweight with no external dependencies
- Simple API with functional options for customization

## Installation

```bash
go get github.com/yourneighborhoodchef/browserforge
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourneighborhoodchef/browserforge/fingerprint"
)

func main() {
    // Create a new fingerprint generator
    generator, err := fingerprint.New()
    if err != nil {
        log.Fatalf("Error creating generator: %v", err)
    }
    
    // Generate a complete fingerprint
    fp, err := generator.Generate()
    if err != nil {
        log.Fatalf("Error generating fingerprint: %v", err)
    }
    
    // Access fingerprint properties
    fmt.Println("User Agent:", fp.UserAgent)
    fmt.Println("OS/CPU:", fp.OSCpu)
    
    // Access headers
    fmt.Println("Accept-Language header:", fp.Headers["Accept-Language"])
}
```

### Generating Headers Only

```go
// If you only need HTTP headers
headers, err := generator.GenerateHeadersOnly()
if err != nil {
    log.Fatalf("Error generating headers: %v", err)
}

fmt.Println("User-Agent:", headers["User-Agent"])
```

### Advanced Usage with Options

```go
// Create a generator with options
generator, err := fingerprint.NewWithOptions(
    fingerprint.WithBrowser("chrome"),
    fingerprint.WithOperatingSystem("windows"),
)
if err != nil {
    log.Fatalf("Error creating generator: %v", err)
}

// Generate a fingerprint with the specified options
fp, err := generator.Generate()
```

## Command Line Tool

BrowserForge also includes a command-line tool:

```bash
# Generate and print HTTP headers as JSON
browserforge headers

# Generate and print a browser fingerprint without headers
browserforge fingerprint

# Generate and print a complete fingerprint with headers
browserforge all
```

## Project Structure

```
browserforge/
├── fingerprint/           # Main package for fingerprint generation
│   ├── options.go         # Configuration options
│   └── fingerprint.go     # Public API
├── internal/              # Implementation details
│   ├── bayesian/          # Bayesian network implementation
│   │   ├── network.go
│   │   └── node.go
│   ├── headers/           # Header generation
│   │   └── generator.go
│   └── data/              # Embedded data resources
├── examples/              # Usage examples
├── cmd/                   # Command-line tools
│   └── browserforge/
│       └── main.go
└── README.md
```

## License

[MIT License](LICENSE)