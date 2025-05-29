package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gabrielnov/tcc-api/internal/dto"
	"github.com/gabrielnov/tcc-api/internal/utils"
)

const DIRECTORY = "./scanned_code"
const MAX_ITERATIONS = 5

type CodeAnalysisService interface {
	Run(req dto.RequestDto) (dto.ResponseDTO, error)
}

type codeAnalysisService struct {
	LlmService    LlmService
	banditService BanditService
	fileManager   FileManager
}

func NewCodeAnalysisService(llmService LlmService, banditService BanditService, fileManager FileManager) CodeAnalysisService {
	return &codeAnalysisService{LlmService: llmService, banditService: banditService, fileManager: fileManager}
}

func (c *codeAnalysisService) Run(req dto.RequestDto) (dto.ResponseDTO, error) {
	filepath, err := c.writeCodeFile(req.Filename, req.Content)
	defer os.Remove(filepath)

	if err != nil {
		log.Printf("error writing code to file: %s", err.Error())
		return dto.ResponseDTO{}, err
	}

	totalIterations := 0
	currentContent := req.Content
	currentContent = utils.FormatPythonCode(currentContent)

	var response dto.ResponseDTO
	response.Success = false

	for ; totalIterations < MAX_ITERATIONS; totalIterations++ {
		banditResults, err := c.banditService.RunAnalysis(filepath)

		if err != nil {
			log.Printf("error running bandit: %s", err.Error())
			return dto.ResponseDTO{}, err
		}

		// zero security issues, analysis stops here
		if len(banditResults) == 0 {

			if totalIterations == 0 {
				return response, nil
			}

			log.Printf("analysis finished with %d iterations", totalIterations)
			response.Success = true
			response.ResultingCode = currentContent
			response.Iterations = totalIterations

			return response, nil
		}

		for k, b := range banditResults {
			log.Printf("fixing vulnerability %s in iteration %d", k, totalIterations)

			lines := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(b.Lines)), ","), "[]")

			prompt := c.LlmService.GeneratePrompt(currentContent, k, b.IssueText, lines)

			currentContent, err = c.LlmService.CallGenAI(prompt)

			currentContent = utils.FormatPythonCode(currentContent)

			if err != nil {
				log.Printf("error calling llm: %s", err.Error())
				return dto.ResponseDTO{}, err
			}

			_, err = c.writeCodeFile(req.Filename, currentContent)
			if err != nil {
				log.Printf("error writing file %s: %s", filepath, err.Error())
				return dto.ResponseDTO{}, err
			}

			log.Printf("fixed vulnerability %s", k)
		}
	}

	log.Printf("incomplete analysis, finished with %d interations", totalIterations)
	response.ResultingCode = currentContent
	response.Iterations = totalIterations
	return response, nil
}

func (c *codeAnalysisService) writeCodeFile(filename, content string) (string, error) {
	if err := os.MkdirAll(DIRECTORY, 0755); err != nil {
		return "", err
	}

	content = utils.FormatPythonCode(content)

	filepath := fmt.Sprintf("%s/%s", DIRECTORY, filename)
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return "", err
	}

	return filepath, nil

}
