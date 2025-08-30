# Contributing to Super MCP Server

Thank you for your interest in contributing to the Super MCP Server! This document provides guidelines for contributing to this project.

## Getting Started

### Prerequisites

- Go 1.24.4 or later
- The `super` command-line tool installed from https://github.com/brimdata/super
- Git

### Setting up the Development Environment

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/your-username/super-mcp-server.git
   cd super-mcp-server
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   go build
   ```

4. Test the server locally:
   ```bash
   ./super-mcp-server
   ```

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker to report bugs or request features
- Before creating a new issue, check if a similar issue already exists
- Provide clear steps to reproduce bugs
- Include relevant system information (Go version, OS, etc.)

### Submitting Changes

1. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and ensure they follow the project conventions:
   - Follow Go conventions and best practices
   - Keep the single-file architecture (`main.go`)
   - Maintain compatibility with the MCP protocol
   - Test your changes thoroughly

3. Commit your changes with a clear commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   ```

4. Push to your fork and create a pull request:
   ```bash
   git push origin feature/your-feature-name
   ```

### Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and concise
- Follow the existing code structure and patterns

### Testing

- Test your changes with various data formats (JSON, CSV, Parquet, etc.)
- Verify compatibility with Claude Desktop integration
- Test error handling and edge cases
- Ensure the server starts and responds correctly to MCP requests

## Development Guidelines

### Architecture Principles

- Maintain the single-file MCP server design
- Keep external dependencies minimal
- Ensure the server remains stateless
- Follow MCP protocol specifications

### Tool Development

When adding or modifying tools:

- Register tools using the established pattern
- Provide clear parameter descriptions
- Handle errors gracefully
- Return structured responses
- Document tool capabilities

### Documentation

- Update README.md if adding new features
- Update CLAUDE.md for development guidance
- Include inline comments for complex logic
- Document any new dependencies or requirements

## Pull Request Process

1. Ensure your code builds without errors
2. Test the server with sample data
3. Update documentation if needed
4. Create a clear pull request description
5. Link any related issues

## Community Guidelines

- Be respectful and constructive in discussions
- Help others learn and contribute
- Follow the project's code of conduct
- Collaborate openly and transparently

## Questions?

If you have questions about contributing, feel free to:
- Open an issue for discussion
- Reach out to the maintainers
- Check existing documentation and issues

Thank you for contributing to Super MCP Server!