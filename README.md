# Lovable Clone - Minimal MVP

An educational project to build a simplified AI code generator that creates single-file web applications from natural language descriptions.

## Project Goals

This is a learning-focused project to understand:
- LLM API integration (OpenAI/Anthropic)
- Code generation and parsing
- Browser sandboxing with iframes
- Full-stack development (Go backend + React frontend)

## Architecture

### Simplified Design (MVP)

```
┌─────────────────────┐
│   React Frontend    │
│                     │
│  - Text input       │
│  - Send button      │
│  - Loading state    │
│  - Code display     │
│  - Iframe preview   │
└──────────┬──────────┘
           │ HTTP POST /api/generate
           │
           ▼
┌─────────────────────┐
│   Go Backend        │
│   (Single Service)  │
│                     │
│  1. Receive prompt  │
│  2. Call LLM API    │
│  3. Parse response  │
│  4. Return code     │
└─────────────────────┘
```

### What's Included (MVP)

✅ Single Go backend service
✅ React frontend with TypeScript
✅ Simple HTTP POST/response (no WebSocket)
✅ LLM API integration (OpenAI/Anthropic - pluggable)
✅ Code parsing (extract HTML/CSS/JS from LLM response)
✅ Live iframe preview

### What's NOT Included (MVP)

❌ Multiple microservices
❌ WebSocket streaming
❌ User authentication
❌ Session persistence
❌ Code iterations/history
❌ Deployment functionality
❌ API Gateway pattern

## API Specification

### Single Endpoint

```
POST /api/generate
Content-Type: application/json

Request:
{
  "prompt": "Create a calculator with dark mode"
}

Response (Success):
{
  "success": true,
  "code": {
    "html": "<!DOCTYPE html>...",
    "css": "body { margin: 0; }...",
    "js": "console.log('hello');"
  }
}

Response (Error):
{
  "success": false,
  "error": "LLM API call failed: timeout"
}
```

## Project Structure

```
lovable-clone/
├── README.md
├── CLAUDE.md (project instructions)
│
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── handlers/
│   │   └── generate.go      # HTTP handler for /api/generate
│   ├── llm/
│   │   ├── interface.go     # LLM provider interface
│   │   ├── openai.go        # OpenAI implementation
│   │   └── anthropic.go     # Anthropic implementation
│   └── parser/
│       └── parser.go        # Extract code from LLM response
│
└── frontend/
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    └── src/
        ├── App.tsx           # Main application
        ├── components/
        │   ├── PromptInput.tsx
        │   ├── CodeViewer.tsx
        │   └── Preview.tsx
        ├── api/
        │   └── client.ts     # API client
        └── types/
            └── index.ts      # TypeScript types
```

## Implementation Plan

### Phase 1: Backend Foundation (2-3 hours)

**Step 1: Setup**
- Initialize Go module
- Create basic HTTP server
- Add CORS middleware

**Step 2: LLM Integration**
- Define provider interface
- Implement OpenAI client
- Handle API calls and errors

**Step 3: Code Parser**
- Parse markdown code blocks
- Extract HTML/CSS/JS
- Handle malformed responses

**Step 4: Generate Handler**
- Wire everything together
- Request validation
- Error handling

### Phase 2: Frontend (2-3 hours)

**Step 1: Setup**
- Initialize Vite + React + TypeScript
- Basic project structure

**Step 2: Components**
- PromptInput: textarea + submit button
- CodeViewer: display generated code (optional)
- Preview: iframe with sandbox attributes

**Step 3: Integration**
- API client
- State management (useState)
- Loading states
- Error handling

### Phase 3: Integration & Testing (1 hour)

- End-to-end testing
- Bug fixes
- Basic styling

**Total Estimate: 5-7 hours**

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- OpenAI API key OR Anthropic API key

### Environment Variables

Create a `.env` file in the `backend/` directory:

```bash
OPENAI_API_KEY=your_key_here
# OR
ANTHROPIC_API_KEY=your_key_here

LLM_PROVIDER=openai  # or "anthropic"
```

### Running Locally

**Backend:**
```bash
cd backend
go run main.go
# Server runs on http://localhost:8080
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# Dev server runs on http://localhost:3000
```

## LLM System Prompt

The backend uses this system prompt to guide code generation:

```
You are a code generator that creates single-file web applications.

Rules:
- Generate complete, working HTML including inline CSS and JavaScript
- Use modern JavaScript (ES6+)
- Make it visually appealing with CSS
- Include all code in one HTML file
- Be creative and functional

Format your response as:
```html
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
      // JS here
    </script>
  </body>
</html>
```
```

## Testing Strategy

**Approach: Build what you need (TDD-light)**

- Write tests for critical paths
- Test LLM parsing logic (multiple fixtures)
- Test error handling
- Manual E2E testing for MVP

## Future Enhancements (Post-MVP)

Once the MVP is working, consider adding:

1. **Streaming responses** - Real-time text generation
2. **Code iterations** - "Make it blue", "Add more features"
3. **Session storage** - Save conversation history
4. **Multi-file projects** - Generate React apps, not just HTML
5. **Microservices architecture** - Split into multiple services
6. **WebSocket real-time** - Bidirectional communication
7. **Deployment** - One-click deploy to Netlify/Vercel
8. **Authentication** - User accounts and saved projects

## Learning Resources

- [OpenAI API Documentation](https://platform.openai.com/docs)
- [Anthropic API Documentation](https://docs.anthropic.com)
- [Go HTTP Server Tutorial](https://go.dev/doc/tutorial/web-service-gin)
- [React + TypeScript](https://react.dev/learn/typescript)
- [Iframe Sandboxing](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/iframe#sandbox)

## Contributing

This is an educational project. Feel free to:
- Experiment with different approaches
- Add features you want to learn
- Share what you learned

## License

MIT - Use for learning purposes
