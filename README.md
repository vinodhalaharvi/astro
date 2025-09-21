# Astro 🚀

A powerful Go AST analyzer that uses segregated interfaces and generics to explore your codebase's structure,
dependencies, and relationships with topological sorting and NoOp implementation generation.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Examples](#examples)
- [Configuration](#configuration)
- [Architecture Deep Dive](#architecture-deep-dive)
- [Contributing](#contributing)
- [License](#license)

## Overview

Astro is a comprehensive Go AST (Abstract Syntax Tree) analysis tool that helps developers understand their codebase
structure through dependency analysis, topological sorting, and automated code generation. Built with Go generics and
segregated interface design principles, Astro provides maximum flexibility and reusability.

### What Makes Astro Special?

- **🧬 Segregated Interface Design**: Each interface has a single, focused responsibility
- **⚡ Generic Architecture**: Type-safe analysis for any Go construct
- **📊 Dependency Analysis**: Understands relationships between types
- **🔄 Topological Sorting**: Shows dependency hierarchy levels
- **🤖 Code Generation**: Automatically generates NoOp implementations
- **🔍 Multi-Type Analysis**: Analyzes structs, interfaces, functions, variables, constants, and imports

## Features

### Core Analysis

- **Struct Analysis**: Fields, methods, embedded types, and dependencies
- **Interface Analysis**: Method signatures, embedded interfaces, and dependencies
- **Function Analysis**: Parameters, returns, receivers, and dependencies
- **Variable & Constant Analysis**: Types, values, and package information
- **Import Analysis**: Package dependencies and aliases

### Advanced Features

- **Topological Sorting**: Dependency-aware ordering (least dependent first)
- **Dependency Visualization**: Shows what each type depends on
- **Level Assignment**: Assigns dependency levels to each construct
- **NoOp Generation**: Creates stub implementations for interfaces
- **Multiple Sort Strategies**: Topological or alphabetical sorting
- **Flexible Output**: Console output and file generation

### Architecture Benefits

- **Type Safety**: Compile-time guarantees through generics
- **Modularity**: Mix and match components as needed
- **Extensibility**: Easy to add new Go construct types
- **Testability**: Each component can be tested independently
- **Performance**: Single AST traversal with multiple analyzers

## Architecture

Astro uses a **segregated interface design** with **Go generics** for maximum flexibility:

```go
// Core segregated interfaces
type NodeVisitor[T any] interface {
    VisitNode(node ast.Node) T
}

type DependencyExtractor[T any] interface {
    ExtractDependencies(item T) []string
}

type ItemSorter[T any] interface {
    SortItems(items []T) []T
}

type CodeGenerator[T any] interface {
    GenerateCode(item T) string
}
```

Each interface has a **single responsibility**, making the system highly modular and testable.

## Installation

```bash
# Clone the repository
git clone https://github.com/vinodhalaharvi/astro.git
cd astro

# Build the binary
go build -o astro main.go

# Or run directly
go run main.go [flags]
```

## Quick Start

```bash
# Analyze all types in current directory
./astro

# Show only structs and interfaces with dependency order
./astro -structs -interfaces

# Generate NoOp implementations for interfaces
./astro -interfaces -noop

# Analyze specific directories
./astro -dirs="./pkg,./cmd,./internal" -all
```

## Usage

### Command Line Flags

| Flag          | Description                            | Default    |
|---------------|----------------------------------------|------------|
| `-dirs`       | Comma-separated directories to analyze | `"."`      |
| `-structs`    | Show struct analysis                   | `false`    |
| `-interfaces` | Show interface analysis                | `false`    |
| `-functions`  | Show function analysis                 | `false`    |
| `-variables`  | Show variable analysis                 | `false`    |
| `-constants`  | Show constant analysis                 | `false`    |
| `-imports`    | Show import analysis                   | `false`    |
| `-all`        | Show all types                         | `false`    |
| `-topo`       | Use topological sorting                | `true`     |
| `-alpha`      | Use alphabetical sorting               | `false`    |
| `-noop`       | Generate NoOp implementations          | `false`    |
| `-noop-dir`   | Directory for NoOp files               | `"./noop"` |

### Basic Usage

```bash
# Analyze current directory (shows all types by default)
./astro

# Analyze specific types
./astro -structs -interfaces

# Use alphabetical sorting instead of topological
./astro -alpha

# Generate NoOp implementations
./astro -interfaces -noop -noop-dir="./generated"
```

### Advanced Usage

```bash
# Comprehensive analysis with code generation
./astro -all -topo -noop -dirs="./pkg,./internal"

# Focus on architecture understanding
./astro -structs -interfaces -functions -topo

# Generate test helpers
./astro -interfaces -noop -noop-dir="./test/mocks"
```

## Examples

### Example Output

```
=== Analyzing file: example.go ===

--- Interfaces (Dependency Order) ---
[Level 0] Interface: Reader (Package: main) at example.go:10:1
  Methods: Read(p []byte) (int, error)

[Level 1] Interface: Writer (Package: main) at example.go:15:1 [depends on: Reader]
  Methods: Write(p []byte) (int, error), Close() error

--- Structs (Dependency Order) ---
[Level 0] Struct: BaseConfig (Package: main) at example.go:25:1
  Fields: Name string, Version int

[Level 1] Struct: ServiceConfig (Package: main) at example.go:30:1 [depends on: BaseConfig]
  Fields: BaseConfig, Port int, Handler Handler
```

### Generated NoOp Implementation

```go
// NoOpReader is a no-op implementation of Reader interface (Level 0)
type NoOpReader struct {
    level int // Dependency level: 0
}

// NewNoOpReader creates a new no-op implementation at the specified level
func NewNoOpReader(level int) *NoOpReader {
    return &NoOpReader{level: level}
}

// GetLevel returns the dependency level of this NoOpReader
func (n *NoOpReader) GetLevel() int {
    return n.level
}

// Read is a no-op implementation (Level 0)
func (n *NoOpReader) Read(p []byte) (int, error) {
    // TODO: Implement Read (Level 0)
    return 0, nil
}
```

## Configuration

### Project Structure

```
your-project/
├── main.go
├── pkg/
│   ├── models/
│   └── services/
├── internal/
│   └── handlers/
└── generated/          # Generated NoOp implementations
    └── noop/
        ├── noop_models_interfaces.go
        └── noop_services_interfaces.go
```

### Recommended Workflow

1. **Architecture Analysis**: Start with `-structs -interfaces -topo` to understand your system's structure
2. **Dependency Review**: Look for circular dependencies or overly complex dependency chains
3. **Test Setup**: Use `-noop` to generate test doubles for interfaces
4. **Continuous Analysis**: Integrate into CI/CD to track architectural changes

## Architecture Deep Dive

### Segregated Interface Design

Astro follows the **Interface Segregation Principle** strictly. Each interface has a single, focused responsibility:

#### Core Interfaces

```go
// Data extraction
type NodeVisitor[T any] interface {
    VisitNode(node ast.Node) T
    }

type ResultCollector[T any] interface {
    CollectResults() []T
    AddResult(item T)
}

// Validation and filtering
type ItemValidator[T any] interface {
    IsValid(item T) bool
}

// Dependency analysis
type DependencyExtractor[T any] interface {
    ExtractDependencies(item T) []string
}

type TypeNameProvider[T any] interface {
    GetTypeName(item T) string
}

// Sorting strategies
type ItemSorter[T any] interface {
    SortItems(items []T) []T
}

type DependencyResolver[T any] interface {
    ResolveDependencies(items []T) []T
}

// Output formatting
type ItemRenderer[T any] interface {
    RenderItem(item T) string
}

// Code generation
type CodeGenerator[T any] interface {
    GenerateCode(item T) string
}
```

### Generic Components

The system uses composition to build complex functionality from simple interfaces:

```go
// Combines multiple responsibilities through composition
type GenericVisitor[T any] struct {
    nodeVisitor NodeVisitor[T]
    collector   ResultCollector[T]
    validator   ItemValidator[T]
}

type AnalysisEngine[T any] struct {
    visitor       *GenericVisitor[T]
    sorter        ItemSorter[T]
    formatter     *GenericFormatter[T]
    codeGenerator *GenericCodeGenerator[T]
}

```

### Benefits of This Design

1. **Single Responsibility**: Each interface does one thing well
2. **Open/Closed Principle**: Easy to extend without modifying existing code
3. **Dependency Inversion**: High-level modules don't depend on low-level modules
4. **Type Safety**: Generics ensure compile-time type checking
5. **Testability**: Each component can be mocked independently

### Adding New Go Constructs

To add analysis for a new Go construct (e.g., type aliases):

1. **Define the domain type**:

```go
type GoTypeAlias struct {
    Name     string
    Package  string
    Target   string
    Position string
    Level    int
}
```

2. **Implement the segregated interfaces**:

```go
type TypeAliasNodeVisitor struct { /* ... */ }
type TypeAliasResultCollector struct { /* ... */ }
type TypeAliasValidator struct { /* ... */ }
// ... etc
```

3. **Compose them into an analyzer**:

```go
    visitor := NewGenericVisitor(
        NewTypeAliasNodeVisitor(fset, pkg),
        NewTypeAliasResultCollector(),
        &TypeAliasValidator{},
    )
```

## Use Cases

### 1. Architecture Review

```bash
# Understand system dependencies
./astro -structs -interfaces -topo -dirs="./internal,./pkg"
```

### 2. Refactoring Preparation

```bash
# Identify tightly coupled components
./astro -all -topo | grep "depends on" | sort
```

### 3. Test Setup

```bash
# Generate test doubles
./astro -interfaces -noop -noop-dir="./test/mocks"
```

### 4. Documentation Generation

```bash
# Export architecture overview
./astro -all -topo > architecture.txt
```

### 5. CI/CD Integration

```bash
# Track architectural changes
./astro -interfaces -topo | diff - previous_architecture.txt
```

## Best Practices

### 1. Regular Analysis

- Run Astro regularly to catch architectural drift
- Integrate into your CI/CD pipeline
- Track dependency complexity over time

### 2. Dependency Management

- Aim for shallow dependency hierarchies
- Watch for circular dependencies
- Consider breaking up complex types with many dependencies

### 3. Interface Design

- Use generated NoOp implementations as starting points
- Keep interfaces focused and cohesive
- Review interface dependencies for design issues

### 4. Code Organization

- Use topological sorting to understand build order
- Organize packages by dependency level
- Consider dependency direction in package structure

## Troubleshooting

### Common Issues

**1. Parse Errors**

```
Error: failed to parse file.go: expected declaration, found 'IDENT'
```

- **Solution**: Ensure Go files are syntactically correct
- **Check**: Run `go vet` and `go build` first

**2. Missing Dependencies**

```
Warning: dependency 'SomeType' not found in current analysis
```

- **Solution**: Include all relevant directories with `-dirs`
- **Check**: Ensure all source files are being analyzed

**3. Circular Dependencies**

```
Warning: circular dependency detected between TypeA and TypeB
```

- **Solution**: This indicates a design issue that should be addressed
- **Check**: Review the dependency chain and consider refactoring

### Performance Considerations

- **Large Codebases**: Use `-dirs` to limit analysis scope
- **Memory Usage**: For very large projects, analyze in chunks
- **Build Time**: Generated NoOp files can increase compilation time

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/vinodhalaharvi/astro.git
cd astro

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o astro main.go
```

### Adding Features

1. **Follow segregated interface design**
2. **Add comprehensive tests**
3. **Update documentation**
4. **Provide examples**

## Roadmap

- [ ] **Web UI**: Browser-based dependency visualization
- [ ] **Graph Export**: DOT/GraphViz output for dependency graphs
- [ ] **Metrics**: Complexity metrics and architectural health scores
- [ ] **Plugin System**: Custom analyzers and generators
- [ ] **Multi-Language**: Support for other languages beyond Go
- [ ] **IDE Integration**: VSCode and GoLand plugins

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Go's `go/ast` package design
- Built with principles from Clean Architecture
- Uses Interface Segregation and Dependency Inversion principles

---

**Happy analyzing! 🚀**

For questions, issues, or contributions, please visit our [GitHub repository](https://github.com/vinodhalaharvi/astro).