package internal

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/EpykLab/secdoc/models"
)

func ParseDirectory(dir string) (*models.Report, error) {
	report := &models.Report{}
	fset := token.NewFileSet()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			controls, requirements, err := ParseFile(fset, path)
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

func ParseFile(fset *token.FileSet, filePath string) ([]models.SecurityControl, []models.Requirement, error) {
	// Parse the Go source file
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var controls []models.SecurityControl
	var requirements []models.Requirement

	// Process all comment groups in the file
	for _, cgroup := range node.Comments {
		text := cgroup.Text()
		pos := fset.Position(cgroup.Pos())

		// Check if this is a security control comment block
		if strings.Contains(text, "@security-control") {
			control := ParseSecurityControl(text, filePath, pos)
			if control != nil {
				controls = append(controls, *control)
			}
		}

		// Check if this is a requirement comment block
		if strings.Contains(text, "@requirement") {
			req := ParseRequirement(text, filePath, pos)
			if req != nil {
				requirements = append(requirements, *req)
			}
		}
	}

	return controls, requirements, nil
}

func ParseSecurityControl(text, filePath string, pos token.Position) *models.SecurityControl {
	control := &models.SecurityControl{
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

func ParseRequirement(text, filePath string, pos token.Position) *models.Requirement {
	req := &models.Requirement{
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
