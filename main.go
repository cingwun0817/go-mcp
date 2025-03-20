package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	calculatorTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Hello World Tool"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Says hello to the name"),
		),
	)

	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["name"].(string)

		return mcp.NewToolResultText("Hello " + name), nil
	})

	// srv := server.NewStdioServer(s)
	// srv.Listen(context.Background(), os.Stdin, os.Stdout)

	srv := server.NewSSEServer(s)
	log.Printf("SSE server listening on localhost:8081\n")
	srv.Start("localhost:8081")

	// if err := server.ServeStdio(s); err != nil {
	// 	fmt.Printf("Server error: %v\n", err)
	// }
}
