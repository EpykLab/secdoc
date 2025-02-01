package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type SecurityControl struct {
	ControlID    string         `json:"control_id"`
	Description  string         `json:"description"`
	References   string         `json:"references"`
	Verification string         `json:"verification"`
	FilePath     string         `json:"file_path"`
	Position     token.Position `json:"position"`
}

type Requirement struct {
	RequirementID string         `json:"requirement_id"`
	Description   string         `json:"description"`
	Verification  string         `json:"verification"`
	Stakeholder   string         `json:"stakeholder"`
	FilePath      string         `json:"file_path"`
	Position      token.Position `json:"position"`
}

type Report struct {
	SecurityControls []SecurityControl `json:"security_controls"`
	Requirements     []Requirement     `json:"requirements"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: secparser <source-dir> [output-file]")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	outputPath := "security-report.json"
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	report, err := parseDirectory(sourcePath)
	if err != nil {
		fmt.Printf("Error parsing directory: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report written to %s\n", outputPath)
}

func parseDirectory(dir string) (*Report, error) {
	report := &Report{}
	fset := token.NewFileSet()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			controls, requirements, err := parseFile(fset, path)
			if err != nil {
				return fmt.Errorf("error parsing %s: %v", path, err)
			}
			report.SecurityControls = append(report.SecurityControls, controls...)
			report.Requirements = append(report.Requirements, requirements...)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return report, nil
}

func parseFile(fset *token.FileSet, filePath string) ([]SecurityControl, []Requirement, error) {
	// Parse the Go source file
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var controls []SecurityControl
	var requirements []Requirement

	// Process all comment groups in the file
	for _, cgroup := range node.Comments {
		text := cgroup.Text()
		pos := fset.Position(cgroup.Pos())

		// Check if this is a security control comment block
		if strings.Contains(text, "@security-control") {
			control := parseSecurityControl(text, filePath, pos)
			if control != nil {
				controls = append(controls, *control)
			}
		}

		// Check if this is a requirement comment block
		if strings.Contains(text, "@requirement") {
			req := parseRequirement(text, filePath, pos)
			if req != nil {
				requirements = append(requirements, *req)
			}
		}
	}

	return controls, requirements, nil
}

func parseSecurityControl(text, filePath string, pos token.Position) *SecurityControl {
	control := &SecurityControl{
		FilePath: filePath,
		Position: pos,
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "@security-control"):
			control.ControlID = strings.TrimSpace(strings.TrimPrefix(line, "@security-control"))
		case strings.HasPrefix(line, "@description:"):
			control.Description = strings.TrimSpace(strings.TrimPrefix(line, "@description:"))
		case strings.HasPrefix(line, "@references:"):
			control.References = strings.TrimSpace(strings.TrimPrefix(line, "@references:"))
		case strings.HasPrefix(line, "@verification:"):
			control.Verification = strings.TrimSpace(strings.TrimPrefix(line, "@verification:"))
		}
	}

	if control.ControlID == "" {
		return nil
	}
	return control
}

func parseRequirement(text, filePath string, pos token.Position) *Requirement {
	req := &Requirement{
		FilePath: filePath,
		Position: pos,
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "@requirement"):
			req.RequirementID = strings.TrimSpace(strings.TrimPrefix(line, "@requirement"))
		case strings.HasPrefix(line, "@description:"):
			req.Description = strings.TrimSpace(strings.TrimPrefix(line, "@description:"))
		case strings.HasPrefix(line, "@verification:"):
			req.Verification = strings.TrimSpace(strings.TrimPrefix(line, "@verification:"))
		case strings.HasPrefix(line, "@stakeholder:"):
			req.Stakeholder = strings.TrimSpace(strings.TrimPrefix(line, "@stakeholder:"))
		}
	}

	if req.RequirementID == "" {
		return nil
	}
	return req
}
