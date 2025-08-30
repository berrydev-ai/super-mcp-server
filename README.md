# Super MCP Server

A Lambda-deployable MCP server with embedded SuperDB for data querying.

## Features

- **Embedded SuperDB**: No external dependencies, Super binary is embedded in the Go executable
- **Lambda Ready**: Runs in AWS Lambda with automatic binary extraction to `/tmp`
- **Local Development**: Also works as a local MCP server for testing with Claude Desktop
- **Three Core Tools**:
  - `query_data`: Execute SuperSQL queries on data files or URLs
  - `list_files`: List data files in directories
  - `analyze_data`: Analyze data structure and types

## Setup

1. **Download Super binary**:
```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

2. **Build for local testing**:
```bash
go mod tidy
go build -o super-mcp-server
```

3. **Build for Lambda**:
```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip lambda-deployment.zip bootstrap
```

## Local Usage

Configure Claude Desktop (`claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "super-data": {
      "command": "/absolute/path/to/super-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

Then restart Claude Desktop and you can ask questions like:
- "List the data files in my current directory"
- "Query my sample.json file and show the first 10 records"
- "Analyze the structure of this CSV file"

## Lambda Deployment

1. Upload `lambda-deployment.zip` to AWS Lambda
2. Set runtime to "Provide your own bootstrap on Amazon Linux 2"
3. Set handler to "bootstrap" 
4. Increase timeout (recommend 30+ seconds for large queries)
5. Increase memory (recommend 512MB+ for better performance)

## Project Structure

```
super-mcp-server/
├── main.go                    # Complete server implementation
├── go.mod                     # Go module dependencies
├── scripts/
│   └── setup.sh              # Downloads Super binary
├── binaries/
│   └── super-linux-amd64     # Embedded Super binary (created by setup.sh)
└── README.md                 # This file
```

## Example Queries

Once connected to Claude Desktop, you can ask:

### Basic Queries
- "What data files do I have in my current directory?"
- "Show me the first 5 records from data.json"
- "Count the rows in my CSV file"

### Advanced Analytics
- "Find the top 10 users by activity in my logs"
- "Group sales data by month and show totals" 
- "Analyze the data types in my JSON file"

### URL Data Sources
- "Query this API endpoint and show me the structure: https://api.example.com/data.json"
- "Download and analyze data from this URL"

## SuperSQL Language

The server uses SuperDB's SuperSQL dialect, which extends SQL with:
- **Pipe syntax**: `SELECT * FROM data.json | WHERE age > 25 | GROUP BY department`
- **JSON native**: No need for special JSON functions
- **Schema-less**: Works with heterogeneous data
- **Rich types**: Supports complex nested data structures

## Technical Notes

- The Super binary is embedded at compile time (~15-20MB)
- In Lambda, the binary is extracted to `/tmp` on first invocation
- Supports data from files, URLs, S3 paths (if Lambda has permissions)
- Works with JSON, CSV, Parquet, NDJSON, and many other formats
- Container reuse means binary extraction only happens once per container lifecycle

## Troubleshooting

### "Super binary not found" error
- Make sure you ran `./scripts/setup.sh` before building
- Check that `binaries/super-linux-amd64` exists and is executable

### Claude Desktop connection issues
- Use absolute paths in the configuration
- Restart Claude Desktop after configuration changes
- Check that the binary is executable: `chmod +x super-mcp-server`

### Lambda timeout errors  
- Increase Lambda timeout for large data processing
- Consider using Lambda layers for very large deployments
- Monitor CloudWatch logs for specific error messages
