# Where I Left Off - Lovable Clone Project

**Date:** 2025-10-20
**Status:** Planning complete, ready to start implementation

---

## What We've Accomplished

### ✅ Completed

1. **Defined project scope** - Minimal MVP with single Go backend + React frontend
2. **Created RFC** - Should be in README.md (need to create)
3. **Architecture decisions:**
   - Single Go backend (no microservices for MVP)
   - HTTP POST/response (no WebSocket streaming)
   - Pluggable LLM interface (OpenAI/Anthropic)
   - Simple iframe preview
   - No auth, no sessions, no persistence

### 📚 Learning Completed

Covered these concepts in detail:
- Go interfaces and why we use them
- Constructor pattern in Go (`NewOpenAIClient()`)
- Why we store `*http.Client` in structs (connection pooling)
- Struct tags for JSON marshaling (`json:"field"`)
- How Go imports work (exported vs unexported)
- Memory layout (heap vs stack vs data segment)
- Dead code elimination in Go linker

---

## Current Status: Ready to Implement Backend

### Next Immediate Steps

#### 1. Create Project Structure

```bash
# From project root
mkdir -p backend/handlers backend/llm backend/parser
cd backend
```

#### 2. Initialize Go Module

```bash
# From backend/ directory
go mod init lovable-clone/backend
```

This creates `go.mod` file.

#### 3. Create LLM Provider Interface

**File:** `backend/llm/interface.go`

```go
package llm

import "context"

// Provider defines the interface that all LLM providers must implement
type Provider interface {
    // GenerateCode takes a user prompt and returns generated code
    GenerateCode(ctx context.Context, prompt string) (string, error)
}
```

#### 4. Implement OpenAI Client

**File:** `backend/llm/openai.go`

```go
package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

// OpenAIClient implements the Provider interface for OpenAI's API
type OpenAIClient struct {
    apiKey     string
    model      string
    httpClient *http.Client
}

// Request structures
type openAIRequest struct {
    Model    string          `json:"model"`
    Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// Response structures
type openAIResponse struct {
    Choices []openAIChoice `json:"choices"`
    Error   *openAIError   `json:"error,omitempty"`
}

type openAIChoice struct {
    Message openAIMessage `json:"message"`
}

type openAIError struct {
    Message string `json:"message"`
    Type    string `json:"type"`
}

// System prompt that guides the LLM
const systemPrompt = `You are a code generator that creates single-file web applications.

Rules:
- Generate complete, working HTML including inline CSS and JavaScript
- Use modern JavaScript (ES6+)
- Make it visually appealing with good CSS styling
- Include all code in one HTML file
- Be creative and functional
- Do not include explanations, only code

You MUST format your response with the code inside a markdown code block like this:

` + "```html" + `
<!DOCTYPE html>
<html>
<head>
    <style>
        /* CSS here */
    </style>
</head>
<body>
    <!-- HTML here -->
    <script>
        // JavaScript here
    </script>
</body>
</html>
` + "```" + `

Important: Only return the code block, no additional text before or after.`

// NewOpenAIClient creates a new OpenAI client with default configuration
func NewOpenAIClient() *OpenAIClient {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        panic("OPENAI_API_KEY environment variable is required")
    }

    return &OpenAIClient{
        apiKey: apiKey,
        model:  "gpt-4",
        httpClient: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

// GenerateCode implements Provider.GenerateCode for OpenAI
func (c *OpenAIClient) GenerateCode(ctx context.Context, prompt string) (string, error) {
    // 1. Build the request payload
    reqBody := openAIRequest{
        Model: c.model,
        Messages: []openAIMessage{
            {
                Role:    "system",
                Content: systemPrompt,
            },
            {
                Role:    "user",
                Content: prompt,
            },
        },
    }

    // 2. Marshal to JSON
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    // 3. Create HTTP request with context
    req, err := http.NewRequestWithContext(
        ctx,
        "POST",
        "https://api.openai.com/v1/chat/completions",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    // 4. Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    // 5. Make the request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to call OpenAI API: %w", err)
    }
    defer resp.Body.Close()

    // 6. Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }

    // 7. Check for HTTP errors
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
    }

    // 8. Parse response
    var openAIResp openAIResponse
    if err := json.Unmarshal(body, &openAIResp); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }

    // 9. Check for API errors
    if openAIResp.Error != nil {
        return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
    }

    // 10. Extract the generated text
    if len(openAIResp.Choices) == 0 {
        return "", fmt.Errorf("no response from OpenAI")
    }

    return openAIResp.Choices[0].Message.Content, nil
}
```

---

## What Comes Next (After OpenAI Client)

### Step 5: Code Parser

**File:** `backend/parser/parser.go`

**Purpose:** Extract HTML/CSS/JS from LLM's markdown response

**What it needs to do:**
1. Find code blocks (```html ... ```)
2. Extract the HTML content
3. Optionally split out inline CSS and JS (or keep as single HTML string)
4. Handle errors (no code block found, malformed HTML)

### Step 6: HTTP Handler

**File:** `backend/handlers/generate.go`

**Purpose:** Handle POST /api/generate endpoint

**Flow:**
1. Parse request JSON (`{"prompt": "..."}`)
2. Call LLM provider
3. Parse code from response
4. Return JSON (`{"success": true, "code": {...}}`)

### Step 7: Main Server

**File:** `backend/main.go`

**Purpose:** HTTP server entry point

**What it needs:**
1. Create HTTP server
2. Register /api/generate handler
3. Add CORS middleware (allow frontend to call it)
4. Listen on :8080

### Step 8: Environment Setup

**File:** `backend/.env.example`

```bash
OPENAI_API_KEY=sk-your-key-here
# OR
ANTHROPIC_API_KEY=your-key-here
```

---

## Testing Plan

Once backend is complete:

```bash
# Terminal 1: Run backend
cd backend
export OPENAI_API_KEY="sk-..."
go run main.go

# Terminal 2: Test with curl
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Create a simple calculator"}'
```

Expected response:
```json
{
  "success": true,
  "code": {
    "html": "<!DOCTYPE html>..."
  }
}
```

---

## Important Commands

```bash
# Check Go version
go version

# Initialize module
go mod init lovable-clone/backend

# Download dependencies (auto-downloads on first build)
go mod download

# Build (checks for errors)
go build

# Run
go run main.go

# Format code
go fmt ./...

# Run tests (when we write them)
go test ./...
```

---

## Files to Create (Checklist)

- [ ] `backend/llm/interface.go`
- [ ] `backend/llm/openai.go`
- [ ] `backend/parser/parser.go`
- [ ] `backend/handlers/generate.go`
- [ ] `backend/main.go`
- [ ] `backend/.env.example`
- [ ] `README.md` (project RFC)

---

## Frontend (Later)

After backend works:

1. Initialize Vite + React + TypeScript
2. Create PromptInput component
3. Create Preview component (iframe)
4. API client to call backend
5. Wire everything together

---

## Key Concepts Learned

### Go Patterns
- Interface-based design (dependency injection)
- Constructor functions for initialization
- Error wrapping with `fmt.Errorf` and `%w`
- Context for cancellation/timeouts
- Struct tags for JSON marshaling

### HTTP Client Best Practices
- Reuse `http.Client` for connection pooling
- Set timeouts to prevent hanging
- Use `defer resp.Body.Close()` to prevent leaks
- Use `http.NewRequestWithContext` for cancellation

### Memory & Performance
- Package-level variables live in data segment
- Linker eliminates unused code (dead code elimination)
- Only pay for what you use

---

## When You Return to Coding

1. **Start here:** Create the directory structure and initialize Go module
2. **Then:** Implement files in this order:
   - interface.go
   - openai.go
   - parser.go
   - generate.go
   - main.go
3. **Test:** Use curl to test the backend
4. **Then:** Move to frontend

---

## Questions to Ask Me When You Resume

- "Should I implement the parser now?"
- "How do I test the OpenAI client?"
- "What should the parser return format be?"
- "How do I add CORS middleware?"

---

## Current Todo Status

- [x] Project structure planned
- [x] Architecture decided
- [x] Go concepts learned
- [ ] Create backend directories
- [ ] Initialize Go module
- [ ] Implement LLM interface
- [ ] Implement OpenAI client
- [ ] Implement parser
- [ ] Implement HTTP handler
- [ ] Implement main server
- [ ] Test with curl
- [ ] Build frontend

---

**Remember:** Start simple, get it working, then iterate!
