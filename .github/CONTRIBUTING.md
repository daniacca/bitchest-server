# Contributing to Bitchest

Thank you for your interest in contributing to Bitchest! Any kind of contributions are welcome!
This document provides guidelines and information for contributors.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style and Standards](#code-style-and-standards)
- [Testing](#testing)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)
- [Feature Requests](#feature-requests)

## Getting Started

Before contributing, please:

1. Read the [README.md](../README.md) to understand the project
2. Familiarize yourself with the [supported commands](../README.md#-supported-commands)
3. Check existing [issues](https://github.com/daniacca/bitchest-server/issues) to avoid duplicates

## Development Setup

### Prerequisites

- **Go 1.24.3 or later** - [Download here](https://golang.org/dl/)
- **Node.js 16+** - For git hooks and commit linting
- **Docker** (optional) - For containerized development

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/daniacca/bitchest-server.git
   cd bitchest-server
   ```

2. **Install dependencies**

   ```bash
   # Install Go dependencies
   go mod download

   # Install Node.js dependencies for git hooks
   npm ci
   ```

3. **Set up git hooks**

   ```bash
   npm run prepare
   ```

4. **Build and test**
   ```bash
   make build-all
   make test
   ```

### Available Make Commands

```bash
# Build targets
make build          # Build server binary
make build-cli      # Build CLI client binary
make build-all      # Build both server and CLI

# Test targets
make test           # Run all tests
make test-verbose   # Run tests with verbose output
make coverage       # Generate coverage report

# Run targets
make run            # Start server on localhost:7463
make run-cli        # Run CLI client
make run-port       # Start server on port 6379
make run-host       # Start server on all interfaces

# Docker targets
make docker-build   # Build Docker image
make docker-run     # Run Docker container

# Help
make help           # Show all available commands
```

## Code Style and Standards

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for code formatting (automatically applied via git hooks)
- Run `go vet` to check for common mistakes
- Keep functions small and focused
- Add comments for exported functions and complex logic

### Project Structure

```
bitchest/
├── cmd/           # Application entry points
│   ├── server/    # Server binary
│   └── cli/       # CLI client binary
├── internal/      # Private application code
├── doc/           # Documentation and assets
├── out/           # Build outputs
└── .github/       # GitHub-specific files
```

### Naming Conventions

- Use descriptive names for variables, functions, and packages
- Follow Go naming conventions (camelCase for variables, PascalCase for exported)
- Use clear, descriptive commit messages. Follows the conventional commit guideline.

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Generate coverage report
make coverage
```

### Test Requirements

- **Coverage**: Maintain at least 60% test coverage overall the codebase (PR will not be merged otherwise)
- **Unit Tests**: Write tests for all new functionality
- **Integration Tests**: Test command interactions and RESP protocol
- **Test Naming**: Use descriptive test names that explain the scenario

### Test Structure

```go
func TestCommandName_Scenario(t *testing.T) {
    // Arrange
    // Act
    // Assert
}
```

## Commit Guidelines

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for commit messages. The format is:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Commit Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools

### Examples

```bash
feat: add EXPIRE command support
fix(server): handle nil responses correctly
docs: update README with new commands
test: add coverage for TTL command
```

### Git Hooks

The project uses Lefthook for automated quality checks:

- **commit-msg**: Validates conventional commit format
- **pre-commit**: Runs `go fmt`, `go vet`, and tests on staged files
- **pre-push**: Runs full test suite and build verification

## Pull Request Process

### Before Submitting

1. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**

   - Write clear, well-documented code
   - Add tests for new functionality
   - Update documentation if needed

3. **Run quality checks**

   ```bash
   make test
   make coverage
   go vet ./...
   go fmt ./...
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

### Pull Request Guidelines

1. **Title**: Use conventional commit format
2. **Description**: Clearly describe the changes and their purpose
3. **Tests**: Ensure all tests pass and coverage is maintained
4. **Documentation**: Update README or docs if needed
5. **Screenshots**: Include screenshots for UI changes (if applicable)

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

1. **Environment**: OS, Go version, Node.js version
2. **Steps to reproduce**: Clear, step-by-step instructions
3. **Expected behavior**: What you expected to happen
4. **Actual behavior**: What actually happened
5. **Additional context**: Logs, screenshots, or other relevant information

### Issue Template

Use the provided issue template when creating new issues. This helps us understand and resolve issues more efficiently.

## Feature Requests

### Before Requesting Features

1. Check if the feature already exists
2. Review existing issues for similar requests
3. Consider if the feature aligns with Bitchest's goals

### Feature Request Guidelines

- **Clear description**: Explain what you want to achieve
- **Use case**: Describe the problem it solves
- **Implementation ideas**: Suggest how it might be implemented
- **Priority**: Indicate if it's a nice-to-have or critical feature

## Getting Help

- **Documentation**: Check the [README.md](../README.md) first
- **Issues**: Search existing [issues](https://github.com/daniacca/bitchest-server/issues)
- **Discussions**: Use GitHub Discussions for questions and ideas

## Code of Conduct

This project is committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and constructive in all interactions.

## License

By contributing to Bitchest, you agree that your contributions will be licensed under the [MIT License](../LICENSE).

---

Thank you for contributing to Bitchest! 🚀
