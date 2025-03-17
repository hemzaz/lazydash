# CLAUDE.md - Guidelines for lazydash Codebase

## Commands
- Build/Install: `go build -o lazydash` or `go install`
- Run all tests: `go test -v ./...`
- Run single test: `go test -v -run TestName` (replace TestName with specific test name)
- Run with coverage: `go test -cover ./...`
- Install from GitHub: `go get -u github.com/hemzaz/lazydash`

## Code Style Guidelines

### Imports
- Standard library imports first, then third-party packages
- Group imports by source (stdlib vs external)

### Naming and Formatting
- Follow Go standard naming: CamelCase for exported, camelCase for unexported
- File naming: feature.go with corresponding feature_test.go
- Use descriptive function and variable names that explain purpose

### Error Handling
- Check errors immediately after function calls
- Use log.Fatal for critical errors that should terminate execution
- Return errors up the call stack when appropriate

### Testing
- Name tests with "Test" prefix followed by function being tested
- Use t.Errorf for test assertions
- Use testutil.go helpers for common test operations

### Documentation
- Add comments for all exported functions and types
- Include usage examples in docstrings where helpful