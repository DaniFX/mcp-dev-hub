package docs

import "fmt"

// Generator produces documentation from repository content.
type Generator struct{}

// NewGenerator creates a new documentation generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// FromReadme extracts and summarizes a README content.
func (g *Generator) FromReadme(content string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("empty README content")
	}
	// TODO: integrate LLM summarization or template-based generation
	return fmt.Sprintf("# Summary\n\n%s", content[:min(len(content), 500)]), nil
}

// FromCode generates a doc stub from source file content.
func (g *Generator) FromCode(filename, content string) (string, error) {
	// TODO: parse comments, exported symbols, package docs
	return fmt.Sprintf("## %s\n\n```go\n%s\n```", filename, content), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
