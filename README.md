# SecDoc Documentation

## Overview

Secdoc is a static analysis tool that extracts security control and requirement documentation from
Go source code comments. It uses Go's Abstract Syntax Tree (AST) to reliably parse and analyze source files, generating
a structured report of security controls and mission requirements.

## Installation

```bash
go install github.com/EpykLab/secdoc@latest
```

Or build from source:

```bash
git clone https://github.com/EpykLab/secdoc.git
cd secdoc
go build
```

## Usage

### Command Line Interface

Basic usage:
```bash
secdoc <source-directory> [output-file]
```

Arguments:
- `source-directory`: Path to the Go source code directory to analyze (required)
- `output-file`: Path for the output JSON report (optional, defaults to "security-report.json")

Example:
```bash
secdoc ./src report.json
```

### Comment Syntax

The parser recognizes two types of special comment blocks:

#### 1. Security Controls

```go
// @security-control AC-3
// @description: Implements role-based access control through middleware
// @references: NIST SP 800-53 Rev 5
// @verification: Unit tests verify middleware blocks unauthorized access

func authMiddleware(next http.Handler) http.Handler {
    // ... implementation
}
```

Fields:
- `@security-control`: (Required) The control identifier
- `@description`: Description of how the code implements the control
- `@references`: Reference to standards or documentation
- `@verification`: How to verify the control is working

#### 2. Mission Requirements

```go
// @requirement REQ-123
// @description: Implements file size validation before upload
// @verification: E2E tests verify large files are rejected
// @stakeholder: Security Team

func validateFileSize(size int64) error {
    // ... implementation
}
```

Fields:
- `@requirement`: (Required) The requirement identifier
- `@description`: Description of how the code implements the requirement
- `@verification`: How to verify the requirement is met
- `@stakeholder`: Who requested/needs this requirement

### Output Format

The tool generates a JSON report with the following structure:

```json
{
  "security_controls": [
    {
      "control_id": "AC-3",
      "description": "Implements role-based access control through middleware",
      "references": "NIST SP 800-53 Rev 5",
      "verification": "Unit tests verify middleware blocks unauthorized access",
      "file_path": "internal/middleware/auth.go",
      "position": {
        "filename": "internal/middleware/auth.go",
        "offset": 1234,
        "line": 42,
        "column": 1
      }
    }
  ],
  "requirements": [
    {
      "requirement_id": "REQ-123",
      "description": "Implements file size validation before upload",
      "verification": "E2E tests verify large files are rejected",
      "stakeholder": "Security Team",
      "file_path": "internal/upload/validation.go",
      "position": {
        "filename": "internal/upload/validation.go",
        "offset": 5678,
        "line": 156,
        "column": 1
      }
    }
  ]
}
```

## Integration Examples

### GitHub Actions

```yaml
name: Security Control Documentation

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  document-controls:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install secdoc
        run: go install github.com/EpykLab/secdoc@latest

      - name: Generate security control report
        run: secdoc ./src security-controls.json

      - name: Upload report artifact
        uses: actions/upload-artifact@v3
        with:
          name: security-controls
          path: security-controls.json
```

### Pre-commit Hook

```bash
#!/bin/bash

# .git/hooks/pre-commit

secdoc . security-controls.json

if [ $? -ne 0 ]; then
    echo "Error: Failed to generate security control documentation"
    exit 1
fi

git add security-controls.json
```

## Development and Extension

### Adding New Comment Types

To add a new type of comment parsing:

1. Define a new struct type for your comment:

```go
type NewCommentType struct {
    ID          string         `json:"id"`
    Description string         `json:"description"`
    FilePath    string         `json:"file_path"`
    Position    token.Position `json:"position"`
}
```

2. Add it to the Report struct:

```go
type Report struct {
    SecurityControls []SecurityControl `json:"security_controls"`
    Requirements    []Requirement     `json:"requirements"`
    NewComments    []NewCommentType  `json:"new_comments"`
}
```

3. Add a parser function:

```go
func parseNewComment(text, filePath string, pos token.Position) *NewCommentType {
    // Implement parsing logic
}
```

4. Update the file parser to recognize the new comment type:

```go
func parseFile(fset *token.FileSet, filePath string) (*Report, error) {
    // ... existing code ...

    if strings.Contains(text, "@new-comment") {
        newComment := parseNewComment(text, filePath, pos)
        if newComment != nil {
            report.NewComments = append(report.NewComments, *newComment)
        }
    }
}
```

### AST Analysis Extensions

The current implementation focuses on comment parsing, but you can extend it to analyze the code itself:

```go
func analyzeImplementation(node ast.Node) {
    ast.Inspect(node, func(n ast.Node) bool {
        switch x := n.(type) {
        case *ast.FuncDecl:
            // Analyze function implementations
            analyzeFunctionSecurity(x)
        case *ast.CallExpr:
            // Track security-critical function calls
            analyzeSecurityCalls(x)
        }
        return true
    })
}
```

## Troubleshooting

### Common Issues

1. **No controls found**: Verify your comment syntax matches the expected format exactly. Comments must start with the exact tag (e.g., `@security-control`).

2. **Parser errors**: If you get parser errors:
    - Verify the target directory contains Go files
    - Check file permissions
    - Ensure Go files are valid and can be compiled

3. **Missing information**: All fields after the initial tag are optional. However, for documentation quality, try to provide all relevant fields.

### Debug Mode

Run with debug logging:
```bash
DEBUG=1 secdoc ./src report.json
```

## Best Practices

1. **Comment Placement**: Place security control comments immediately above the relevant code implementation.

2. **Control IDs**: Use standardized control IDs from your compliance framework (e.g., NIST 800-53, ISO 27001).

3. **Verification**: Always include specific, testable verification steps.

4. **Updates**: Keep security control documentation updated when modifying security-critical code.

5. **Reviews**: Include security control documentation review in your code review process.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on:
- Submitting bug reports
- Creating pull requests
- Development setup
- Testing requirements

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
