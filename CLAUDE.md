# Super MCP Server Development Guide

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a Model Context Protocol (MCP) server implementation in Go that provides SuperSQL querying capabilities. It acts as a bridge between Claude Desktop and the Super data processing tool, allowing natural language queries to be executed against various data formats.

## Core Architecture
- **Single-file MCP server** (`main.go`) implementing two primary tools:
  - `query`: Executes SuperSQL queries via the external `super` command
  - `list_files`: Discovers data files in directories
- **Tool registration pattern**: Uses the MCP Go SDK's tool registration system
- **External dependency**: Requires the `super` command-line tool to be installed and in PATH

## Code Style Guidelines
- **Imports**: Standard library first, third-party second, internal packages last
- **Error Handling**: Use if err != nil with meaningful error returns/wrapping
- **Testing**: Use testify for assertions, mock interfaces when needed
- **Naming**: 
  - Variables/fields: camelCase
  - Functions/methods: CamelCase
  - Packages: single lowercase word
- **Commits and Pull Requests**: 
  - Commits should be atomic and have a descriptive commit message
  - Pull requests should be made against the `main` branch
  - Don't mention claude in the commit message

### Planning Phase
1. Use `gh issue view` to get the issue details
2. Understand the problem described in the issue
3. Ask clarifying questions if necessary
4. Understand the prior art for this issue
   - Search the scratchpads for previous thoughts related to the issue
   - Search PRs to see if you can find history on this issue
   - Search the codebase for relevant files
5. Think harder about how to break the issue down into a series of small, manageable tasks
6. Document your plan in a new scratchpad
   - Include the issue name in the filename
   - Include a link to the issue in the scratchpad

### Implementation Phase
7. Create a new branch for the issue
8. Solve the issue in small, manageable steps, according to your plan
9. Commit your changes after each step. Only add the new or updated files related to the current task
10. Never commit with `-a` flag to add all the files in the repo

### Testing Phase
- Write Go tests to describe the expected behavior of your code
- Run the full test suite to ensure you haven't broken anything
- If the tests are failing, fix them
- Ensure that all tests are passing before moving on to the next step

### Deployment Phase
- Open a PR and request a review

## GitHub Actions Integration

### Issue Research Workflow
When mentioned in a new issue with `@claude`:

1. **Data Gathering Phase**:
   - Read and understand the issue description thoroughly
   - Search the codebase for related functionality using Grep/Glob tools
   - Review previous issues and PRs for similar problems or implementations
   - Check documentation and existing tests for context
   - Identify relevant files, functions, and components

2. **Research Analysis**:
   - Analyze the current implementation state
   - Identify potential root causes or implementation approaches  
   - Document any existing technical debt or dependencies
   - Note any architectural considerations or constraints

3. **Summary Generation**:
   - Create a comprehensive research summary including:
     - **Current State**: What exists today related to this issue
     - **Key Findings**: Important discoveries from codebase analysis
     - **Related Components**: Files/functions/APIs involved
     - **Implementation Notes**: Technical considerations and approach recommendations
     - **Dependencies**: Any blockers or prerequisites identified
   - Append the summary as a comment to the original issue
   - Use clear headings and bullet points for readability

### Code Review Integration
- Automatic code review triggered on PR creation/updates
- Focus on Go best practices, security, and OnScript API patterns
- Reference existing codebase patterns and conventions

## Security Guidelines
- **API Key Protection**: Never hardcode secrets, always use GitHub Secrets
- **Workflow Security**: Follow principle of least privilege for permissions
- **Secret Validation**: All workflows include security validation steps
- **Monitoring**: Weekly automated security audits and manual reviews
- See `docs/security-guidelines.md` for comprehensive security practices

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