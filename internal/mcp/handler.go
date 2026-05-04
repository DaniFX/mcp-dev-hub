package mcp

// Handler routes MCP JSON-RPC 2.0 requests to the appropriate tool.
type Handler struct {
	tools map[string]ToolFunc
}

// ToolFunc is the signature for all MCP tool implementations.
type ToolFunc func(params map[string]interface{}) (interface{}, error)

// NewHandler creates a new MCP handler with registered tools.
func NewHandler() *Handler {
	return &Handler{
		tools: make(map[string]ToolFunc),
	}
}

// Register adds a tool to the handler.
func (h *Handler) Register(name string, fn ToolFunc) {
	h.tools[name] = fn
}

// Call dispatches a tool call by name.
func (h *Handler) Call(name string, params map[string]interface{}) (interface{}, error) {
	fn, ok := h.tools[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return fn(params)
}
