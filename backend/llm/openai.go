package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIClient struct {
	httpClient *http.Client
	apiKey     string
	model      string
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

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

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	if apiKey == "" {
		panic("API Key must be provided")
	}

	if model == "" {
		model = "gpt-4"
	}

	return &OpenAIClient{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		apiKey: apiKey,
		model:  model,
	}
}

func (o *OpenAIClient) GenerateCode(ctx context.Context, prompt string) (string, error) {
	request := openAIRequest{
		Model: o.model,
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

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	bufferJson := bytes.NewReader(jsonData)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bufferJson)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error (status %d) %s", resp.StatusCode, string(body))
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}
