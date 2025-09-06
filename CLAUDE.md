# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a toolkit service - a web-based collection of utility tools built with Go (Gin framework) and modern frontend technologies. The application provides 10 different utility tools through a single web interface with a modern, animated design inspired by Magic UI.

## Architecture

### Backend Structure
- **Main Entry Point**: `main.go` - Contains the Gin web server setup, route handlers, and tokenizer integration
- **Web Framework**: Gin v1.7.4 with Go 1.14+
- **Tokenizer Integration**: Custom tokenizer package for text processing functionality
- **Static File Serving**: All frontend assets served from `resources/static/`

### Frontend Structure
- **Main Application**: `resources/static/index.html` - Landing page with tool navigation
- **Individual Tools**: Each tool has its own subdirectory in `resources/static/`
- **Shared Assets**: Common CSS and JavaScript in `resources/static/css/` and `resources/static/js/`

### Tools Available
1. HTML to Image Converter (`/html2img`)
2. JSON Formatter (`/json-formatter`)
3. Base64 Encoder/Decoder (`/base64-encoder`)
4. Regex Tester (`/regex-tester`)
5. URL Encoder/Decoder (`/url-encoder`)
6. Hash Calculator (`/hash-calculator`)
7. Timestamp Converter (`/timestamp-converter`)
8. UUID Generator (`/uuid-generator`)
9. Color Picker (`/color-picker`)
10. Tokenizer Analyzer (`/tokenizer`)

## Development Commands

### Running the Application
```bash
# Install dependencies
go mod tidy

# Run development server
go run main.go

# Run with custom port
PORT=3000 go run main.go
```

### Building
```bash
# Build executable (creates ./app)
./build.sh

# Or build manually
go build -o ./app main.go
```

### Testing
```bash
# Run tests (if any)
go test ./...

# Run with verbose output
go test -v ./...
```

## Key Implementation Details

### Main.go Structure
- **Global Variables**: `tools` slice contains tool metadata, `globalTokenizer` for text processing
- **Initialization**: `ConfigRuntime()` sets CPU cores, `initTokenizer()` loads tokenizer config
- **Route Handlers**: Each tool has a dedicated handler function serving static HTML
- **Tokenizer API**: `/api/tokenizer` endpoint supports tokenize, encode, and decode operations

### Tokenizer Integration
- **Configuration**: Loads from `tokenizer/tokenizer.json`
- **Fallback**: Creates basic tokenizer if config loading fails
- **API Modes**: Supports tokenize (analysis), encode (text to token IDs), decode (token IDs to text)

### Static File Organization
- Each tool is self-contained in its own directory
- Main index.html provides navigation and search functionality
- Shared assets are in common css/ and js/ directories

## Environment Configuration

### Environment Variables
- `PORT`: Server port (default: 8080)

### File Locations
- Static files: `resources/static/`
- Tokenizer config: `tokenizer/tokenizer.json`
- Build output: `./app`

## Important Notes

- The build script references `rooms.go`, `routes.go`, and `stats.go` files that don't exist in the current codebase
- The application was refactored from a chat room service to a toolkit service
- The tokenizer integration is a key feature with its own API endpoint
- All tools are currently client-side HTML/JavaScript applications served statically