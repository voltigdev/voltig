# ⚡️ Voltig: System Package Manager Automation and Task Runner

A modern, cross-platform Go-based CLI tool for managing system packages with a unified interface. Voltig supports parallel operations, robust configuration, and developer-friendly features.

---

## ✨ Features
- 🍺 **Homebrew Automation**: Install, remove, update, and check status of Homebrew packages via CLI(Will be extended to support other package managers in the future)
- 📝 **YAML Configuration**: Manage multiple packages and custom commands in `voltig.yml`.
- 🔍 **Automatic Config Discovery**: Voltig finds `voltig.yml` in your current or parent directories.
- 🧩 **Custom Commands**: Define your own build, test, or install scripts in YAML.
- 🚀 **Easy Installation**: One-liner install script or manual build.

---

## 🚀 Installation

### 🍎 MacOS Install

```sh
brew install voltig
```

### Other option:

1. **Clone the repository**
   ```sh
   git clone <repo-url>
   cd voltig-cli
   ```
2. **Install dependencies**
   ```sh
   go mod tidy
   ```
3. **Build and install Voltig**
   ```sh
   go build -o voltig .
   mv ./voltig $HOME/go/bin/
   # Make sure $HOME/go/bin is in your PATH
   ```
---

## 🗂️ Configuration Discovery
- Voltig will automatically find `voltig.yml` in your current directory or any parent directory.
- You can run Voltig commands from any subfolder in your project tree.

---

## 📄 voltig.yml Configuration Guide

The `voltig.yml` file defines packages and custom commands for the Voltig CLI. Below are the allowed formats and best practices for configuration.

### Packages Section

Define packages to be installed and managed. Each package entry can use either a single name (as a string) or multiple names (as an array of strings).

#### Allowed Syntax

```yaml
packages:
  # Single package name (string)
  - name: "python3"
    version: latest
    manager: brew

  # Multiple package names (array)
  - name: ["node", "rust"]
    version: latest
    manager: brew
```

**Field Reference:**
- `name`:  
  - *Type*: string or array of strings  
  - *Description*: The name(s) of the package(s). If multiple, use YAML array syntax.
- `version`:  
  - *Type*: string  
  - *Description*: The version to install (e.g., `latest`, `1.0.0`). Optional.
- `manager`:  
  - *Type*: string  
  - *Description*: The package manager to use (e.g., `brew`).
- `optional`:  
  - *Type*: boolean  
  - *Description*: If true, package is optional. Default is false.
- `dependencies`:  
  - *Type*: array of strings  
  - *Description*: List of package dependencies. Optional.

#### Examples

```yaml
packages:
  # Single package
  - name: "python3"
    version: "3.11"
    manager: brew

  # Multiple packages grouped
  - name: ["node", "rust"]
    version: latest
    manager: brew

  # With dependencies and optional flag
  - name: ["git"]
    manager: brew
    optional: true
    dependencies: ["curl", "openssl"]
```

### Commands Section

Define custom commands to run with Voltig. Each command can have a summary, a shell command, a script, and optional arguments.

```yaml
commands:
  build:
    summary: Build the project
    command: go build -o voltig .
  test:
    summary: Run tests
    command: go test ./...
  setup-dev:
    summary: Run setup script
    script: ./scripts/setup-dev.sh
  deploy:
    summary: Deploy with arguments
    script: ./scripts/deploy.sh
    args:
      - "--prod"
```

**Field Reference:**
- `summary`: Short description of the command.
- `command`: Shell command to execute.
- `script`: Path to a script file to run.
- `args`: List of arguments to pass to the script.

---

### Formatting Tips

- Indentation should be two spaces.
- For arrays, use YAML array syntax:  
  `name: ["node", "rust"]`
- For single values, use quotes for clarity:  
  `name: "python3"`
- Comments start with `#`.

---

**Best Practices:**
- Group related packages using the array syntax for clarity.
- Use descriptive summaries for commands.
- Keep your configuration DRY and organized.

---


## ⚡️ Example Usage

**Install all packages from voltig.yml:**
```sh
voltig install
```

**Remove all packages from voltig.yml:**
```sh
voltig remove
```

**Update all packages from voltig.yml:**
```sh
voltig update
```

**Check status of all packages:**
```sh
voltig status
```

**Install a single package:**
```sh
voltig install zig
```

**Remove a single package:**
```sh
voltig remove gleam
```

---

## 🗃️ Directory Structure
```
├── 📦 cmd/            # Command implementations
├── 📝 config/         # Configuration loader & search logic
├── 🔒 internal/       # Internal packages
│   ├── 📦 manager/    # Package manager interfaces & implementations
│   ├── 📄 models/     # Data models
├── 🚀 main.go         # Entry point
└── 📄 go.mod          # Go module definition
```

## License
MIT
