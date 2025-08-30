# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a Model Context Protocol (MCP) server implementation in Go that provides SuperSQL querying capabilities. It acts as a bridge between Claude Desktop and the Super data processing tool, allowing natural language queries to be executed against various data formats.

## Core Architecture
- **Single-file MCP server** (`main.go`) implementing two primary tools:
  - `query`: Executes SuperSQL queries via the external `super` command
  - `list_files`: Discovers data files in directories
- **Tool registration pattern**: Uses the MCP Go SDK's tool registration system
- **External dependency**: Requires the `super` command-line tool to be installed and in PATH

## Development Commands

### Build
```bash
go mod tidy
go build
```

### Run locally
```bash
./super-mcp-server
```

## Dependencies
- **Super command-line tool**: Must be installed from https://github.com/brimdata/super
- **MCP Go SDK**: `github.com/modelcontextprotocol/go-sdk v0.3.0`
- **Go 1.24.4** or compatible

## Integration
This server is designed to integrate with Claude Desktop through the MCP protocol. The configuration requires adding the built binary path to `claude_desktop_config.json` under the `mcpServers` section.

## Data Format Support
The server supports querying multiple data formats through Super:
- JSON (.json)
- CSV (.csv)
- TSV (.tsv)
- Parquet (.parquet)
- NDJSON (.ndjson)
- Log files (.log)

## Tool Parameters
- **query tool**: Accepts `query` (SuperSQL), optional `data_path` and `format`
- **list_files tool**: Accepts optional `directory` parameter (defaults to current directory)