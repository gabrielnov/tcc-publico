package service

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type SecurityIssue struct {
	TestId     string `json:"test_id"`
	Filename   string `json:"filename"`
	IssueText  string `json:"issue_text"`
	LineNumber int    `json:"line_number"`
	LineRange  []int  `json:"line_range"`
}

type BanditResults struct {
	Results []SecurityIssue `json:"results"`
}

type BanditCompiledResults struct {
	IssueText string `json:"issue_text"`
	Filename  string `json:"filename"`
	Lines     []int  `json:"lines"`
}

type BanditService interface {
	RunAnalysis(filePath string) (map[string]BanditCompiledResults, error)
	compileResults(results []SecurityIssue) map[string]BanditCompiledResults
}

type banditService struct{}

func NewBanditService() BanditService {
	return &banditService{}
}

func (b *banditService) RunAnalysis(filePath string) (map[string]BanditCompiledResults, error) {
	log.Printf("running bandit analysis for file %s", filePath)
	cmd := exec.Command("bandit", "-f", "json", filePath)
	output, _ := cmd.CombinedOutput()

	lines := strings.Split(string(output), "\n")
	if len(lines) > 4 {
		output = []byte(strings.Join(lines[4:], "\n"))
	}

	var banditResults BanditResults
	err := json.Unmarshal(output, &banditResults)

	if err != nil {
		log.Print("Error unmarshaling bandit response, ", err.Error())
		return nil, err
	}

	return b.compileResults(banditResults.Results), nil
}

func (b *banditService) compileResults(results []SecurityIssue) map[string]BanditCompiledResults {
	resultsMap := make(map[string]BanditCompiledResults)

	for _, v := range results {
		_, ok := resultsMap[v.TestId]
		if ok {
			temp := resultsMap[v.TestId]
			temp.Lines = append(temp.Lines, v.LineNumber)
			resultsMap[v.TestId] = temp
		} else {
			resultsMap[v.TestId] = BanditCompiledResults{
				IssueText: v.IssueText,
				Filename:  v.Filename,
				Lines:     []int{v.LineNumber},
			}
		}
	}
	return resultsMap
}
