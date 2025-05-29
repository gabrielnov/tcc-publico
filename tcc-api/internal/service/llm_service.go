package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const OLLAMA_PORT = 11434
const OLLAMA_HOST = "http://localhost"
const OLLAMA_MODEL = "deepseek-coder-v2:16b"

const OLLAMA_PROMPT = `

You will act as a Python code vulnerability fixer. You will receive a test ID from the Bandit static analyzer, ` +
	`a brief description of the vulnerability, a list of lines where the vulnerability is present, and the code containing the vulnerability. ` +
	`You must fix this vulnerability in all the indicated lines. Do not fix vulnerabilities that were not mentioned in the prompt.` +
	`Write comments in code to explain what was done in each and every line that you change. ` +
	`Do not remove code that is required for the proper functioning of the original program. ` +
	`You will return the code and only the code, without any extra information. Any additional information will invalidate your response.

The input is: 
Bandit test id:
%s
Vulnerability description: 
%s 
Lines with vulnerability: 
%s
Vulnerable code snippet: 
%s`

type GenerateRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

type LlmResponse struct {
}

type LlmResponseLine struct {
	Response string `json:"response"`
}

type LlmService interface {
	CallGenAI(prompt string) (string, error)
	GeneratePrompt(code, testId, issueText, lines string) string
}

type llmService struct{}

func NewLlmService() LlmService {
	return &llmService{}
}

func (c *llmService) CallGenAI(prompt string) (string, error) {

	log.Printf("calling llm")
	body := GenerateRequest{Prompt: prompt, Model: OLLAMA_MODEL}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	resp, err := http.Post(fmt.Sprintf("%s:%d/api/generate", OLLAMA_HOST, OLLAMA_PORT), "application/json", payloadBuf)

	if err != nil {
		log.Printf("error sending request to llm: %s", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var completeResponse string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var llmResp LlmResponseLine
		if err := json.Unmarshal(scanner.Bytes(), &llmResp); err != nil {
			log.Printf("error unmarshalling llm response: %s", err.Error())
			return "", err
		}
		completeResponse += llmResp.Response
	}
	return completeResponse, nil

}

func (c *llmService) GeneratePrompt(code, testId, issueText, lines string) string {
	return fmt.Sprintf(OLLAMA_PROMPT, testId, issueText, lines, code)
}
