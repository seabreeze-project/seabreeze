# CLAUDE.md - Seabreeze Development Guide for AI Assistants

**Last Updated**: 2025-11-13
**Project Version**: Early Development (Pre-1.0)
**Go Version**: 1.21+

## Table of Contents

1. [Project Overview](#project-overview)
2. [Codebase Architecture](#codebase-architecture)
3. [Key Components](#key-components)
4. [Development Workflows](#development-workflows)
5. [Code Conventions](#code-conventions)
6. [Configuration System](#configuration-system)
7. [Error Handling Patterns](#error-handling-patterns)
8. [Testing Guidelines](#testing-guidelines)
9. [Common Tasks](#common-tasks)
10. [Important Considerations](#important-considerations)

---

## Project Overview

### What is Seabreeze?

Seabreeze is a lightweight container orchestration tool designed as a "poor man's Kubernetes" for managing Docker Compose-based web applications. It provides centralized management for multiple compose projects with built-in scripting capabilities and an interactive shell interface.

### Core Capabilities

- **Multi-Project Management**: Organize Docker Compose projects across multiple "bases" (directories)
- **Script Execution Engine**: Run commands in containers, services, or on the host with shell variable expansion
- **Interactive Shell**: Command-line shell interface for Seabreeze operations
- **Docker API Proxy**: Restricted access to Docker API endpoints (partial implementation)

### Design Philosophy

- **Simplicity over Complexity**: Avoid Kubernetes complexity while providing essential orchestration
- **Configuration-Driven**: YAML-based project manifests (`seabreeze.yml`)
- **Docker Compose First**: Works with existing compose files, no migration needed
- **Developer-Friendly**: Interactive shell, clear commands, familiar patterns

---

## Codebase Architecture

### Directory Structure

```
seabreeze/
├── main.go                 # Entry point → cmd.Execute()
├── go.mod / go.sum         # Dependencies (Go 1.21+)
├── bin/                    # Build output directory
├── cmd/                    # CLI commands (Cobra-based)
├── config/                 # Configuration models
├── console/                # Output formatting helpers
├── core/                   # Application core (Docker client, config)
├── projects/               # Project & container management
├── proxy/                  # Docker API reverse proxy
├── scripts/                # Script execution engine
├── shell/                  # Interactive shell mode
└── util/                   # Utilities (paths, hashing)
```

### Architectural Patterns

**Pattern**: Modular layer architecture with dependency injection
- **Core Layer** (`core/`): Central coordinator, provides Docker client and config
- **Domain Layer** (`projects/`, `scripts/`, `proxy/`): Business logic
- **CLI Layer** (`cmd/`): User interface via Cobra commands
- **Support Layer** (`util/`, `console/`, `config/`): Cross-cutting concerns

**Dependency Flow**:
```
main.go → cmd.Execute()
  ↓
cmd/* commands → core.Core instance
  ↓
core.Core → projects.*, scripts.*, proxy.*
  ↓
Docker API / filesystem
```

---

## Key Components

### 1. Core Package (`/core/`)

**Purpose**: Central coordinator providing Docker client and configuration access.

**Location**: `/home/user/seabreeze/core/`

**Key Types**:
```go
type Core struct {
    ConfigBasePath string           // Platform-specific config directory
    client         *client.Client   // Lazy-loaded Docker client
    config         *config.Configuration
}
```

**Key Methods**:
- `Initialize()`: Loads configuration from YAML + environment variables
- `Client()`: Returns Docker API client (lazy-loaded)
- `Config()`: Returns application configuration
- `LoadProjectRepository()`: Creates ProjectRepository for project discovery

**Platform-Specific Paths**:
- **Linux**: `/etc/seabreeze`
- **macOS**: `/Library/Application Support/Seabreeze`
- **Windows**: `%PROGRAMDATA%\Seabreeze`

**Docker Setup**:
- Sets `DOCKER_API_VERSION=1.45` to prevent version mismatch
- Auto-detects socket: Unix (`/var/run/docker.sock`) or Windows (`npipe:////./pipe/docker_engine`)

**File References**:
- Configuration logic: `core/services.go:22-78`
- Docker client: `core/services.go:80-108`

### 2. Projects Package (`/projects/`)

**Purpose**: Manages Docker Compose projects, containers, and project metadata.

**Location**: `/home/user/seabreeze/projects/`

**Key Types**:

```go
// Represents a single project directory
type Project struct {
    ID          string  // SHA1 hash of path
    Name        string  // Directory name
    Path        string  // Absolute path
    BaseName    string  // Base identifier (@main, @staging)
    Base        string  // Base directory path
    Description string  // From seabreeze.yml
}

// Manages project discovery across bases
type ProjectRepository struct {
    bases map[string]string  // Named bases (@main → /var/projects)
}

// Project manifest (seabreeze.yml)
type ProjectMetadata struct {
    Description string
    Scripts     map[string]ScriptConfig
}

// Parsed docker-compose.yml
type ComposeConfig struct {
    Services map[string]ServiceConfig
}

// Container within a project
type Container struct {
    ID      string
    Name    string
    Project *Project
    Service string
    Replica int
    State   string
}
```

**Key Methods**:

**ProjectRepository**:
- `FindByIdentifier(id string)`: Resolve `@base/project`, `/absolute/path`, or project name
- `FindAll()`: Discover all projects across all bases
- `Create(base, name string)`: Create new project from template

**Project**:
- `HasManifest()`: Check if `seabreeze.yml` exists
- `HasComposeFile()`: Check if `docker-compose.yml` exists
- `Metadata()`: Parse `seabreeze.yml`
- `ComposeConfig()`: Parse `docker-compose.yml`
- `Services()`: List compose services
- `Container(service, replica)`: Resolve container by service + replica number
- `Status(service)`: Get running container count

**Container Naming Convention**:
- **Static**: Use `container_name` from compose file
- **Dynamic**: `{project-dir}-{service}-{replica}` (e.g., `myapp-web-1`)

**Base Resolution**:
- `@main/myproject` → Look in "main" base
- `@staging/app` → Look in "staging" base
- `/absolute/path` → Direct path
- `myproject` → Search all bases for directory name

**File References**:
- Project model: `projects/project.go:17-54`
- Repository: `projects/repository.go:14-163`
- Container resolution: `projects/container.go:11-121`

### 3. Scripts Package (`/scripts/`)

**Purpose**: Execute commands in containers, services, or on the host with shell expansion.

**Location**: `/home/user/seabreeze/scripts/`

**Key Types**:

```go
type ScriptConfig struct {
    Command     string            // Shell command to execute
    Target      string            // "host", "container name", "service name"
    WorkDir     string            // Working directory
    User        string            // Unix user (host scripts only)
    Environment map[string]string // Environment variables
}

type Executor struct {
    project *projects.Project
}
```

**Execution Targets**:

| Target Type | Format | Example | Execution |
|------------|--------|---------|-----------|
| Host | `host` | `host` | Runs on host OS |
| Container | `container <name>` | `container web-1` | `docker exec` |
| Service | `service <name>` | `service web` | `docker compose exec` |

**Shell Variable Expansion**:

Uses `mvdan.cc/sh/v3` for POSIX-compliant shell expansion:

```yaml
scripts:
  deploy:
    command: "echo Deploying $1 to $2 with args: $@"
    target: "host"
```

**Special Variables**:
- `$@` - All arguments as separate quoted strings
- `$*` - All arguments as single string
- `$#` - Argument count
- `$1`, `$2`, ... - Positional arguments
- **Protected**: These cannot be overridden by environment variables

**Platform Support**:
- **Unix/Linux**: Full support including user switching
- **Windows**: Host scripts cannot switch users (returns error)

**File Organization**:
- Core logic: `scripts/executor.go`
- Unix-specific: `scripts/executor_unix.go` (`//go:build !windows`)
- Windows-specific: `scripts/executor_windows.go` (`//go:build windows`)

**File References**:
- Executor: `scripts/executor.go:17-182`
- Shell expansion: `scripts/executor.go:140-166`

### 4. Shell Package (`/shell/`)

**Purpose**: Interactive shell mode with command completion and aliases.

**Location**: `/home/user/seabreeze/shell/`

**Key Types**:

```go
type Shell struct {
    commands    CommandsMap         // Native shell commands
    aliases     CommandsMap         // Command aliases
    dispatchers []CommandDispatcher // Fallback handlers
    Completer   Completer          // Autocomplete
}

type Command interface {
    Name() string
    Aliases() []string
    Run(args []string) error
}

type VirtualCommand interface {
    Command
    IsVirtual() bool  // Display-only, delegates execution
}
```

**Command Registration**:

```go
sh := shell.New()
sh.AddCommand(myCommand)           // Real shell command
sh.AddVirtualCommand(cobraCmd)     // Display in help, execute via dispatcher
sh.AddAlias("ls", "list")          // Alias support
```

**Dispatcher Pattern**:

When command not found, dispatchers are tried in order:
1. Native shell commands
2. Aliases
3. Registered dispatchers (e.g., Cobra command fallback)

**Autocomplete System**:

Context-aware completion based on:
- Command name matching
- Argument position
- Custom completers per command

**File References**:
- Shell core: `shell/shell.go:18-204`
- Command interface: `shell/command.go:7-25`

### 5. Proxy Package (`/proxy/`)

**Purpose**: Docker API reverse proxy with endpoint filtering (based on `docker-proxy-acl`).

**Location**: `/home/user/seabreeze/proxy/`

**Supported Endpoints**:
- `/containers/*` - Container operations
- `/images/*` - Image management
- `/volumes/*` - Volume operations
- `/networks/*` - Network management
- `/services/*` - Swarm services
- `/tasks/*` - Swarm tasks
- `/events` - Event stream (streaming)
- `/version`, `/info`, `/ping` - System info

**Architecture**:
- **Gorilla Mux** for routing with API version prefix support
- **Unix Socket Proxying**: `source.sock` → `target.sock`
- **Permission Management**: Socket ownership and mode (e.g., `0660`)

**Configuration**:

```yaml
reverse_proxy:
  network_name: seabreeze-network
```

**File References**:
- Proxy server: `proxy/proxy.go:20-187`

### 6. Config Package (`/config/`)

**Purpose**: Configuration models and validation.

**Location**: `/home/user/seabreeze/config/`

**Configuration Structure**:

```go
type Configuration struct {
    Bases        BasesMap           // Named project directories
    ReverseProxy ReverseProxyConfig
}

type BasesMap map[string]string
// Example: {"main": "/var/projects", "staging": "/opt/staging"}

type ReverseProxyConfig struct {
    NetworkName string `mapstructure:"network_name"`
}
```

**Validation Rules**:
- `main` base is **required**
- Viper automatically unmarshals from YAML

**File References**:
- Models: `config/config.go:8-30`

### 7. Console Package (`/console/`)

**Purpose**: Helper for consistent command output formatting.

**Location**: `/home/user/seabreeze/console/`

**Usage Pattern**:

```go
func runCommand(cmd *cobra.Command, args []string) {
    h := console.NewHelper(cmd)

    h.Status("Starting operation...")       // [seabreeze cmd] Starting operation...
    h.Success("Operation complete")         // [seabreeze cmd] Operation complete
    h.Warning("Potential issue")            // [seabreeze cmd] Potential issue
    h.Error("Something failed")             // [seabreeze cmd] Something failed
    h.Fatal(err)                            // Print error and exit

    table := h.Table()
    table.AddRow("Key", "Value")
    h.Print(table)
}
```

**File References**:
- Helper implementation: `console/helper.go`

### 8. Util Package (`/util/`)

**Purpose**: Utility functions for common operations.

**Location**: `/home/user/seabreeze/util/`

**Key Utilities**:

```go
// Path wrapper for filesystem operations
type Path struct { /* private */ }

p := util.NewPath("/base")
p.Join("sub", "file.txt")                    // /base/sub/file.txt
p.Locate("config.yml", "config.yaml")        // Find first existing file
p.Exists()                                   // bool
p.MustExist()                                // panic if not exists
p.IsDir()                                    // bool

// Hashing
util.Sha1(input string) string               // SHA1 hash (project IDs)

// Environment
util.GetEditor() string                      // $EDITOR or fallback
```

**File References**:
- Path utilities: `util/path.go`
- Misc utilities: `util/util.go`

---

## Development Workflows

### Building the Project

```bash
# Linux/macOS
go build -o bin/seabreeze

# Windows
go build -o bin\seabreeze.exe

# Install locally
go install
```

**No Build System**: Project uses simple `go build` without Makefile or build scripts.

### Running During Development

```bash
# Direct execution
go run main.go [command]

# Using built binary
./bin/seabreeze [command]
```

### Adding a New Command

**Step 1**: Create command file in `/cmd/`

```go
// cmd/mycommand.go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/seabreeze-project/seabreeze/console"
)

var myCmd = &cobra.Command{
    Use:   "my-command",
    Short: "Description",
    Args:  cobra.NoArgs,
    Run:   runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)  // Register with root
}

func runMyCommand(cmd *cobra.Command, args []string) {
    h := console.NewHelper(cmd)
    h.Status("Running...")
    // Implementation
}
```

**Step 2**: Add flags if needed

```go
func init() {
    myCmd.Flags().StringP("output", "o", "", "Output file")
    rootCmd.AddCommand(myCmd)
}
```

**Step 3**: Test manually (no automated tests currently)

```bash
go run main.go my-command
```

### Adding Project Features

**Example**: Add new field to project metadata

**Step 1**: Update `ProjectMetadata` in `projects/metadata.go`

```go
type ProjectMetadata struct {
    Description string
    Scripts     map[string]ScriptConfig
    NewField    string `yaml:"new_field"`  // Add this
}
```

**Step 2**: Update template in `projects/template.go` if needed

**Step 3**: Test with real project

```bash
# Create test project
./bin/seabreeze project create @main/test-project

# Edit seabreeze.yml
vim /var/seabreeze/projects/test-project/seabreeze.yml

# Test loading
./bin/seabreeze project info test-project
```

### Platform-Specific Code

Use build tags for OS-specific implementations:

```go
//go:build !windows
// +build !windows

package mypkg

// Unix-specific implementation
```

```go
//go:build windows
// +build windows

package mypkg

// Windows-specific implementation
```

**Existing Pattern**: See `scripts/executor_unix.go` and `scripts/executor_windows.go`

---

## Code Conventions

### Go Style Guidelines

**Follow Standard Go Conventions**:
- Use `gofmt` formatting (spaces not tabs, except for Go requires tabs for indentation)
- Exported names start with capital letter
- Interface names: single-method interfaces end in `-er` (e.g., `Runner`, `Executor`)
- Error variables: prefix with `Err` (e.g., `ErrNotFound`)

### Package Organization

**One Concern Per Package**:
- `projects/` - Only project-related logic
- `scripts/` - Only script execution
- `cmd/` - Only CLI commands

**Avoid Circular Dependencies**:
- `core` can import domain packages (`projects`, `scripts`)
- Domain packages should NOT import `core`
- Use dependency injection via method parameters

### Naming Conventions

**Files**:
- `repository.go` - Multiple repository pattern implementations
- `metadata.go` - Data models/structures
- `executor.go` - Main logic
- `executor_unix.go` - Platform-specific (with build tags)
- `template.go` - Templates

**Types**:
- `Project` - Single entity
- `ProjectRepository` - Collection/manager
- `ProjectMetadata` - Data structure
- `ScriptConfig` - Configuration model

**Methods**:
- `Find()`, `FindByID()`, `FindAll()` - Retrieval
- `Create()`, `Update()`, `Delete()` - Mutations
- `Has...()` - Boolean checks (e.g., `HasManifest()`)
- `Load...()` - Load and return data
- `Must...()` - Panics on failure (e.g., `MustExist()`)

### Error Handling

**Pattern 1**: Early return with Fatal (commands only)

```go
func runCommand(cmd *cobra.Command, args []string) {
    h := console.NewHelper(cmd)

    project, err := repo.FindByIdentifier(args[0])
    if err != nil {
        h.Fatal(err)  // Prints to stderr and exits
    }
    // Continue...
}
```

**Pattern 2**: Error wrapping (library code)

```go
func (r *Repository) FindByID(id string) (*Project, error) {
    project, found := r.projects[id]
    if !found {
        return nil, fmt.Errorf("project %q not found", id)
    }
    return project, nil
}
```

**Pattern 3**: Sentinel errors

```go
var ErrManifestNotFound = errors.New("project manifest not found")

// Later...
if errors.Is(err, ErrManifestNotFound) {
    // Handle specific error
}
```

**Pattern 4**: Panic for developer errors

```go
func (c *Core) Config() *config.Configuration {
    if c.config == nil {
        panic("accessing config before it has been loaded")
    }
    return c.config
}
```

**Avoid**: Global error handlers, structured logging frameworks (keep it simple)

### File Layout

**Standard Order**:
1. Package declaration
2. Imports (stdlib, external, internal)
3. Constants
4. Type definitions
5. Constructor functions (`New...()`)
6. Methods (grouped by receiver)
7. Helper functions

**Example**:
```go
package projects

import (
    "fmt"
    "os"
)

const DefaultReplicas = 1

type Project struct {
    Name string
}

func NewProject(name string) *Project {
    return &Project{Name: name}
}

func (p *Project) String() string {
    return p.Name
}
```

### Documentation

**Exported Symbols**:
- Add doc comments for all exported types, functions, methods
- Start with the symbol name: `// Project represents...`

**TODOs**:
- Use `// TODO: description` for future work
- Include context about why it's needed

**Examples from codebase**:
```go
// TODO: Use label based filtering instead of the container name.
// TODO: History support
```

---

## Configuration System

### Application Configuration

**File**: `{ConfigBasePath}/config.yml`

**Override**: `seabreeze --config /path/to/config.yml`

**Schema**:

```yaml
# config.yml
bases:
  main: /var/seabreeze/projects      # Required
  staging: /opt/staging              # Optional
  production: /srv/production        # Optional

reverse_proxy:
  network_name: seabreeze-network
```

**Environment Variables**:

Prefix: `SEABREEZE_`

```bash
export SEABREEZE_BASES_MAIN=/custom/path
export SEABREEZE_REVERSE_PROXY_NETWORK_NAME=custom-network
```

**Viper Binding**: Automatic mapping from env vars to config structure

### Project Configuration

**File**: `{project-dir}/seabreeze.yml`

**Schema**:

```yaml
# seabreeze.yml
description: "My web application"

scripts:
  # Script on host
  build:
    command: "npm run build"
    target: "host"
    workdir: "/home/user/myproject"

  # Script in specific container
  db-migrate:
    command: "php artisan migrate"
    target: "container myapp-web-1"
    workdir: "/var/www"
    environment:
      DB_HOST: localhost

  # Script in service (any replica)
  logs:
    command: "tail -f /var/log/app.log"
    target: "service web"
    user: "www-data"
```

**Discovery**:
- Searches for `seabreeze.yml` or `seabreeze.yaml`
- Must be in project root alongside `docker-compose.yml`

### Docker Compose Integration

**File**: `{project-dir}/docker-compose.yml`

Seabreeze parses compose files to:
- List available services
- Determine container names and replica counts
- Execute commands in containers

**Not Modified**: Seabreeze reads compose files but never writes to them.

---

## Error Handling Patterns

### Common Error Scenarios

**1. Configuration Not Found**

```go
// core/services.go:69
if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
        // Config file not found, use defaults
        return &config.Configuration{}, nil
    }
    return nil, err  // Other error
}
```

**2. Required Base Missing**

```go
// config/config.go:24
func (c *Configuration) Validate() error {
    if _, ok := c.Bases["main"]; !ok {
        return errors.New("configuration requires 'main' base to be set")
    }
    return nil
}
```

**3. Project Not Found**

```go
// projects/repository.go
func (r *Repository) FindByIdentifier(id string) (*Project, error) {
    // Try various resolution strategies...
    return nil, fmt.Errorf("project %q not found", id)
}
```

**4. Manifest Missing**

```go
// projects/metadata.go:25
var ErrManifestNotFound = errors.New("project manifest not found")

func (p *Project) HasManifest() (bool, error) {
    _, err := p.manifestPath()
    if errors.Is(err, ErrManifestNotFound) {
        return false, nil  // Not an error, just doesn't exist
    }
    if err != nil {
        return false, err  // Real error
    }
    return true, nil
}
```

**5. Platform Limitations**

```go
// scripts/executor_windows.go:11
func executeOnHost(..., user string) error {
    if user != "" {
        return fmt.Errorf("running scripts on the host as specific user is currently not supported on Windows")
    }
    // Continue...
}
```

### Error Handling Decision Tree

```
Is this a CLI command?
├─ Yes → Use console.Helper.Fatal(err)
└─ No → Return error to caller

Is this expected behavior?
├─ Yes → Return nil or specific sentinel error
└─ No → Return wrapped error with context

Is this a developer error (misuse of API)?
├─ Yes → Panic with descriptive message
└─ No → Return error

Can the error be recovered?
├─ Yes → Return error for caller to handle
└─ No → Use Fatal() in commands
```

---

## Testing Guidelines

### Current State

**No Automated Tests**: The project is in early development phase without test coverage.

**Manual Testing**:
1. Build the project: `go build -o bin/seabreeze`
2. Test commands manually with real Docker environments
3. Verify output and behavior

### When Adding Tests (Future)

**Recommended Structure**:

```
projects/
├── project.go
├── project_test.go          # Unit tests
├── repository.go
├── repository_test.go
└── testdata/                # Test fixtures
    ├── sample-project/
    │   ├── seabreeze.yml
    │   └── docker-compose.yml
    └── invalid-project/
```

**Testing Strategies**:

**1. Unit Tests** (when added):
```go
// projects/project_test.go
func TestProject_HasManifest(t *testing.T) {
    // Test logic
}
```

**2. Integration Tests** (when added):
- Use `testdata/` directory for fixtures
- May require Docker daemon for full tests
- Consider using Docker-in-Docker or test containers

**3. Table-Driven Tests** (Go best practice):
```go
func TestSomething(t *testing.T) {
    tests := []struct{
        name string
        input string
        want string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

### Manual Testing Checklist

**Before committing changes**:

- [ ] `go build` succeeds
- [ ] `go fmt ./...` applied
- [ ] No obvious runtime panics
- [ ] Test with real Docker environment if Docker-related
- [ ] Test on target platform (Linux/macOS/Windows if applicable)

---

## Common Tasks

### Task 1: Add a New Project Script Target

**Goal**: Support scripts running in `task` containers (one-off containers).

**Steps**:

1. **Update `ScriptConfig`** in `scripts/executor.go`:
```go
type ScriptConfig struct {
    Command     string
    Target      string  // Now supports: "host", "container X", "service X", "task X"
    // ... other fields
}
```

2. **Add target parsing** in `scripts/executor.go`:
```go
func (e *Executor) Execute(script ScriptConfig, args []string) error {
    parts := strings.Fields(script.Target)
    targetType := parts[0]

    switch targetType {
    case "task":
        return e.executeInTask(parts[1], script, args)
    // ... existing cases
    }
}
```

3. **Implement `executeInTask`**:
```go
func (e *Executor) executeInTask(task string, script ScriptConfig, args []string) error {
    // Use docker compose run instead of exec
    // Similar to executeInService but with 'run --rm'
}
```

4. **Test manually**:
```yaml
# seabreeze.yml
scripts:
  seed:
    command: "php artisan db:seed"
    target: "task seeder"
```

### Task 2: Add Label-Based Container Filtering

**Context**: Currently uses name-based container discovery (see TODO in code).

**Goal**: Use Docker labels for more reliable container identification.

**Steps**:

1. **Update compose integration** to add labels:
   - Consider `com.seabreeze.project=myapp`
   - `com.seabreeze.service=web`

2. **Modify `Project.Status()`** in `projects/project.go`:
```go
func (p *Project) Status(service string) (int, error) {
    filters := filters.NewArgs()
    filters.Add("label", fmt.Sprintf("com.seabreeze.project=%s", p.Name))
    filters.Add("label", fmt.Sprintf("com.seabreeze.service=%s", service))
    filters.Add("status", "running")

    containers, err := p.client.ContainerList(ctx, container.ListOptions{
        Filters: filters,
    })
    // ...
}
```

3. **Document label schema** in CLAUDE.md

4. **Update README** with label requirements

### Task 3: Add Command Validation

**Goal**: Validate commands before execution.

**Steps**:

1. **Create validation package** `validate/`:
```go
// validate/command.go
package validate

func CommandArgs(cmd *cobra.Command, expected int) error {
    if len(os.Args) < expected + 1 {
        return fmt.Errorf("expected %d arguments", expected)
    }
    return nil
}
```

2. **Use in commands**:
```go
func runCommand(cmd *cobra.Command, args []string) {
    h := console.NewHelper(cmd)

    if err := validate.CommandArgs(cmd, 1); err != nil {
        h.Fatal(err)
    }
    // ...
}
```

### Task 4: Improve Error Messages

**Goal**: Add more context to error messages.

**Pattern**:

```go
// Before
return errors.New("not found")

// After
return fmt.Errorf("project %q not found in bases: %v", id, availableBases)
```

**Checklist**:
- [ ] Include identifiers (project name, service name)
- [ ] Suggest solutions where possible
- [ ] List available options when relevant

### Task 5: Add Configuration Validation

**Goal**: Validate configuration on load.

**Steps**:

1. **Extend `Configuration.Validate()`** in `config/config.go`:
```go
func (c *Configuration) Validate() error {
    // Check main base exists
    if _, ok := c.Bases["main"]; !ok {
        return errors.New("configuration requires 'main' base to be set")
    }

    // Validate base paths exist
    for name, path := range c.Bases {
        if !util.NewPath(path).IsDir() {
            return fmt.Errorf("base %q path %q is not a directory", name, path)
        }
    }

    return nil
}
```

2. **Call validation** in `core/services.go` after loading config

---

## Important Considerations

### 1. Development Status

**Early Stage**:
- No tests yet
- APIs may change
- Some features partially implemented
- Marked as "WIP" in README

**Implications for AI Assistants**:
- Don't assume backward compatibility
- Breaking changes are acceptable
- Focus on correct implementation over API stability
- Document design decisions in code comments

### 2. Docker API Versioning

**Critical**: Always set `DOCKER_API_VERSION=1.45`

**Location**: `core/services.go:86`

```go
if _, ok := os.LookupEnv("DOCKER_API_VERSION"); !ok {
    os.Setenv("DOCKER_API_VERSION", "1.45")
}
```

**Why**: Prevents client/server version mismatch errors when Docker updates.

### 3. Platform Differences

**Unix vs Windows**:

| Feature | Unix/Linux | Windows | Implementation |
|---------|------------|---------|----------------|
| Docker Socket | `/var/run/docker.sock` | `npipe:////./pipe/docker_engine` | `core/services.go:95` |
| Config Path | `/etc/seabreeze` | `%PROGRAMDATA%\Seabreeze` | `core/core.go:19` |
| User Switching | Supported via `syscall.Credential` | Not supported | `scripts/executor_*.go` |
| Path Separator | `/` | `\` | Use `filepath.Join()` |

**Always Test Platform-Specific Code**:
- Use build tags: `//go:build !windows`
- Provide fallback implementations
- Document platform limitations

### 4. Container Naming

**Two Patterns**:

**Static** (from compose file):
```yaml
services:
  web:
    container_name: my-web-container  # Fixed name
```

**Dynamic** (generated by Docker Compose):
```yaml
services:
  web:
    deploy:
      replicas: 3  # Creates: myapp-web-1, myapp-web-2, myapp-web-3
```

**Resolution Logic**: `projects/container.go:34-121`

**Implication**: Container name resolution is complex, consider label-based approach (TODO in code).

### 5. Shell Variable Expansion

**Uses Real Shell Parser**: `mvdan.cc/sh/v3`

**Implications**:
- Full POSIX shell syntax supported
- Proper quoting required for arguments with spaces
- `$@` expands correctly (each arg quoted separately)
- Environment variables expanded at execution time

**Security Consideration**:
- User input in `command` field is executed in shell
- Sanitize or validate user-provided commands
- Consider command injection risks

**Example Risk**:
```yaml
scripts:
  dangerous:
    command: "echo $USER_INPUT"  # If USER_INPUT = "foo; rm -rf /"
    target: "host"
```

**Mitigation**: Document that `command` should be trusted input only.

### 6. Configuration Priority

**Order** (lowest to highest):
1. Default values in code
2. `config.yml` file
3. Environment variables (`SEABREEZE_*`)
4. Command-line flags

**Managed by Viper**: Automatic handling of priority.

### 7. Project Identification

**Projects Identified By**:
1. **ID**: SHA1 hash of absolute path (`util.Sha1(path)`)
2. **Name**: Directory name
3. **Identifier**: `@base/name`, `/path`, or `name`

**Implications**:
- Moving project changes its ID (path-based)
- Same directory name in different bases = different projects
- ID used for internal tracking, name for user display

### 8. Docker Compose Integration

**Current Approach**:
- Parses `docker-compose.yml` for metadata
- Shells out to `docker compose` and `docker` commands
- Does not use Compose Go library directly

**TODO in Code**: Use Docker API instead of shelling out.

**Implication**:
- Requires `docker compose` CLI installed
- Cannot work in environments without Docker CLI
- Future: Use Compose Go library for programmatic control

### 9. Error vs. Missing State

**Pattern**: Distinguish between errors and expected missing state.

```go
// GOOD: Returns (false, nil) when manifest doesn't exist
func (p *Project) HasManifest() (bool, error) {
    _, err := p.manifestPath()
    if errors.Is(err, ErrManifestNotFound) {
        return false, nil  // Not an error
    }
    return err == nil, err
}

// BAD: Would return error for missing manifest
func (p *Project) HasManifest() error {
    _, err := p.manifestPath()
    return err  // Missing manifest = error (wrong)
}
```

**Guideline**: Missing optional data should not be an error.

### 10. Template System

**Project Creation**: Uses Go templates

**Location**: `projects/template.go`

**Available in Templates**:
- `{{.Project}}` - Project instance
- `{{.Config}}` - Application configuration

**Example**:
```go
const manifestTemplate = `description: "{{.Project.Name}}"
scripts:
  example:
    command: "echo Hello from {{.Project.Name}}"
    target: "host"
`
```

**Extending Templates**:
1. Add template constants in `template.go`
2. Update `createFromTemplate()` function
3. Create template with proper Go template syntax

---

## Quick Reference

### Key File Paths

| Component | Path | Purpose |
|-----------|------|---------|
| Entry point | `/home/user/seabreeze/main.go` | Application entry |
| Core logic | `/home/user/seabreeze/core/` | Central coordinator |
| CLI commands | `/home/user/seabreeze/cmd/` | Cobra commands |
| Project mgmt | `/home/user/seabreeze/projects/` | Project/container logic |
| Script exec | `/home/user/seabreeze/scripts/` | Script execution |
| Shell mode | `/home/user/seabreeze/shell/` | Interactive shell |
| Proxy | `/home/user/seabreeze/proxy/` | Docker API proxy |
| Config models | `/home/user/seabreeze/config/` | Configuration |
| Console | `/home/user/seabreeze/console/` | Output helpers |
| Utilities | `/home/user/seabreeze/util/` | Shared utils |

### Important Line References

| Feature | File | Line Range |
|---------|------|------------|
| Docker client setup | `core/services.go` | 80-108 |
| Config loading | `core/services.go` | 22-78 |
| Project resolution | `projects/repository.go` | 46-163 |
| Container naming | `projects/container.go` | 34-121 |
| Script execution | `scripts/executor.go` | 48-182 |
| Shell expansion | `scripts/executor.go` | 140-166 |
| Command registration | `cmd/root.go` | - |
| Project metadata | `projects/metadata.go` | 13-43 |

### Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `SEABREEZE_BASES_MAIN` | Main project base | (required) |
| `DOCKER_HOST` | Docker socket | Auto-detected |
| `DOCKER_API_VERSION` | Docker API version | `1.45` |
| `EDITOR` | Text editor | Platform default |

### Build Commands

```bash
# Development
go run main.go [command]

# Build
go build -o bin/seabreeze

# Install globally
go install

# Format
go fmt ./...

# Vet
go vet ./...
```

### Common Patterns

**Creating a Console Helper**:
```go
h := console.NewHelper(cmd)
h.Status("Working...")
h.Success("Done!")
```

**Loading Configuration**:
```go
core := core.New()
core.Initialize()
config := core.Config()
```

**Finding Projects**:
```go
repo := core.LoadProjectRepository()
project, err := repo.FindByIdentifier("@main/myproject")
```

**Executing Scripts**:
```go
executor := scripts.NewExecutor(project)
err := executor.Execute(scriptConfig, args)
```

**Path Operations**:
```go
p := util.NewPath("/base")
exists := p.Join("file.txt").Exists()
```

---

## Contributing Guidelines for AI Assistants

### Code Quality Standards

1. **Follow Existing Patterns**: Study similar code before adding new features
2. **Maintain Consistency**: Use same error handling, naming, structure as existing code
3. **Document Public APIs**: Add doc comments for exported symbols
4. **Add TODOs**: Note future improvements with `// TODO: description`
5. **Platform Awareness**: Use build tags for OS-specific code
6. **Keep It Simple**: Avoid over-engineering in this early-stage project

### Pull Request Checklist

Before suggesting changes:

- [ ] Code builds successfully (`go build`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No new warnings from `go vet ./...`
- [ ] Follows existing code conventions
- [ ] Includes doc comments for exported symbols
- [ ] Platform-specific code uses build tags
- [ ] Error handling follows established patterns
- [ ] Console output uses `console.Helper` in commands
- [ ] New configuration fields added to schema documentation
- [ ] TODOs added for incomplete features

### Communication

**When Explaining Code**:
- Reference file paths and line numbers
- Use examples from the actual codebase
- Explain the "why" not just the "what"
- Note any platform-specific considerations

**When Suggesting Changes**:
- Explain the problem being solved
- Show before/after code
- Note any breaking changes
- Suggest testing approach

---

## Revision History

| Date | Changes | Author |
|------|---------|--------|
| 2025-11-13 | Initial CLAUDE.md creation | AI Assistant |

---

## Additional Resources

- **Repository**: https://github.com/seabreeze-project/seabreeze
- **Issue Tracker**: https://github.com/seabreeze-project/seabreeze/issues
- **License**: MIT (see LICENSE.md)
- **Go Version**: 1.21+ (see go.mod)

**Dependencies Documentation**:
- Cobra: https://github.com/spf13/cobra
- Viper: https://github.com/spf13/viper
- Docker SDK: https://pkg.go.dev/github.com/docker/docker
- Shell Parser: https://pkg.go.dev/mvdan.cc/sh/v3

---

*This document is maintained for AI assistants working on the Seabreeze codebase. Keep it updated as the project evolves.*
