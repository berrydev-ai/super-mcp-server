package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Embed the Super binary
//
//go:embed binaries/super-linux-amd64
var superBinaryData []byte

var (
	superPath string
	setupOnce sync.Once
	setupErr  error
)

// Tool argument types
type QueryArgs struct {
	Query    string `json:"query" jsonschema:"required,description=SuperSQL query to execute"`
	DataPath string `json:"data_path,omitempty" jsonschema:"description=Path to data file or URL"`
	Format   string `json:"format,omitempty" jsonschema:"description=Output format (json, csv, table)"`
}

type ListFilesArgs struct {
	Directory string `json:"directory,omitempty" jsonschema:"description=Directory to list (default: /tmp)"`
}

type AnalyzeArgs struct {
	DataPath string `json:"data_path" jsonschema:"required,description=Path to data file or URL"`
	Sample   int    `json:"sample,omitempty" jsonschema:"description=Number of records to sample (default: 100)"`
}

// Setup Super binary (runs once per Lambda container)
func setupSuperBinary() error {
	setupOnce.Do(func() {
		// Use /tmp directory (writable in Lambda)
		superPath = "/tmp/super"

		// Check if already exists (container reuse)
		if _, err := os.Stat(superPath); err == nil {
			return
		}

		// Write embedded binary to /tmp
		if err := os.WriteFile(superPath, superBinaryData, 0755); err != nil {
			setupErr = fmt.Errorf("failed to write Super binary: %v", err)
			return
		}

		log.Printf("Super binary extracted to %s", superPath)
	})
	return setupErr
}

func handleQuery(ctx context.Context, req *mcp.CallToolRequest, args QueryArgs) (*mcp.CallToolResult, error) {
	if err := setupSuperBinary(); err != nil {
		return errorResult(fmt.Sprintf("Failed to setup Super: %v", err)), nil
	}

	// Build super command
	cmdArgs := []string{}

	// Set output format
	format := args.Format
	if format == "" {
		format = "json"
	}
	cmdArgs = append(cmdArgs, "-f", format)

	// Add query if provided
	if args.Query != "" {
		cmdArgs = append(cmdArgs, "-c", args.Query)
	}

	// Add data path if provided
	if args.DataPath != "" {
		cmdArgs = append(cmdArgs, args.DataPath)
	}

	// Execute super command
	cmd := exec.CommandContext(ctx, superPath, cmdArgs...)

	// Set working directory to /tmp for any temporary files
	cmd.Dir = "/tmp"

	output, err := cmd.CombinedOutput()

	if err != nil {
		return errorResult(fmt.Sprintf("Super command failed: %v\nOutput: %s", err, string(output))), nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(output)},
		},
	}, nil
}

func handleListFiles(ctx context.Context, req *mcp.CallToolRequest, args ListFilesArgs) (*mcp.CallToolResult, error) {
	dir := args.Directory
	if dir == "" {
		dir = "/tmp" // Default to /tmp in Lambda
	}

	// Find data files
	dataExts := []string{".json", ".csv", ".tsv", ".parquet", ".ndjson", ".log", ".jsonl"}
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		for _, ext := range dataExts {
			if strings.HasSuffix(strings.ToLower(path), ext) {
				files = append(files, path)
				break
			}
		}
		return nil
	})

	if err != nil {
		return errorResult(fmt.Sprintf("Error listing files: %v", err)), nil
	}

	result, _ := json.MarshalIndent(map[string]interface{}{
		"directory": dir,
		"files":     files,
		"count":     len(files),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(result)},
		},
	}, nil
}

func handleAnalyze(ctx context.Context, req *mcp.CallToolRequest, args AnalyzeArgs) (*mcp.CallToolResult, error) {
	if err := setupSuperBinary(); err != nil {
		return errorResult(fmt.Sprintf("Failed to setup Super: %v", err)), nil
	}

	sample := args.Sample
	if sample <= 0 {
		sample = 100
	}

	// Build analysis query
	query := fmt.Sprintf("head %d | typeof(this) | sort | uniq -c | sort -r count", sample)

	cmd := exec.CommandContext(ctx, superPath, "-f", "json", "-c", query, args.DataPath)
	cmd.Dir = "/tmp"

	output, err := cmd.CombinedOutput()

	if err != nil {
		return errorResult(fmt.Sprintf("Analysis failed: %v\nOutput: %s", err, string(output))), nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(output)},
		},
	}, nil
}

func errorResult(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: message},
		},
		IsError: &[]bool{true}[0],
	}
}

// Lambda handler
func handleLambda(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {
	// Initialize MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "super-data-server",
			Version: "1.0.0",
		},
		nil,
	)

	// Register tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "query_data",
		Description: "Execute SuperSQL queries on data files or URLs",
	}, handleQuery)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_files",
		Description: "List data files in a directory",
	}, handleListFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "analyze_data",
		Description: "Analyze data structure and types",
	}, handleAnalyze)

	// For Lambda, we need to handle the MCP protocol differently
	// This is a simplified version - in production you'd need proper MCP over HTTP

	return map[string]interface{}{
		"statusCode": 200,
		"body": map[string]interface{}{
			"message": "Super MCP Server initialized",
			"tools":   []string{"query_data", "list_files", "analyze_data"},
		},
	}, nil
}

// Main function - detects if running in Lambda or locally
func main() {
	// Check if running in AWS Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Starting Super MCP Server in Lambda mode")
		lambda.Start(handleLambda)
		return
	}

	// Local development mode - run as MCP server
	log.Println("Starting Super MCP Server in local mode")

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "super-data-server",
			Version: "1.0.0",
		},
		nil,
	)

	// Register tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "query_data",
		Description: "Execute SuperSQL queries on data files or URLs",
	}, handleQuery)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_files",
		Description: "List data files in a directory",
	}, handleListFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "analyze_data",
		Description: "Analyze data structure and types",
	}, handleAnalyze)

	// Run server with stdio transport for local use
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
