# Works On My Machine (womm)

**womm** is a CLI tool that codifies your project's environment requirements into a portable config file and verifies them programmatically — eliminating the "it works on my machine" problem.

Given a `.womm` configuration file, it checks that required tools are installed at the correct versions, environment variables are set, required files exist, and the current platform is supported. It can also install missing dependencies and run setup commands.

## Installation

```sh
go install womm@latest
```

Or build from source:

```sh
cd cli
go build -o womm .
```

## Quick Start

```sh
# Generate a starter .womm file
womm init

# Verify your environment
womm verify

# Check, install missing deps, and run setup
womm
```

## Configuration

Create a `.womm` file in your project root with the following sections:

### `dependencies`

CLI tools that must be installed, with minimum version requirements and optional package-manager-specific install overrides.

```yaml
dependencies:
  node:
    version: ">=18.0"
    brew: node@18
    apt: nodejs
    choco: nodejs-lts

  git:
    version: ">=2.30"

  docker:
    version: ">=20.10"
    brew: --cask docker
```

Supported package managers: `brew`, `apt`, `yum`, `choco`, `scoop`, `winget`.

### `env`

Environment variables that should be set. The value is a hint shown if the variable is missing — it is not written to the environment.

```yaml
env:
  API_URL: http://localhost:3000
  DATABASE_URL: postgres://localhost:5432/app
```

### `files`

Files or directories that must exist in the project root.

```yaml
files:
  - .env
  - config.json
  - certs/
```

### `commands`

Named groups of shell commands. womm asks for confirmation before each group and runs steps sequentially. If a step fails, the remaining steps are skipped.

```yaml
commands:
  setup:
    - npm install
    - npm run build

  seed:
    - node scripts/seed.js
```

### `platforms`

Allowed operating systems. womm exits early if the current OS is not in this list.

```yaml
platforms:
  windows: true
  linux: true
  macos: true
```

Platforms map to `runtime.GOOS` values: `windows`, `linux`, `macos`.

## Usage

| Command | Description |
|---|---|
| `womm` | Check environment, prompt to install missing tools, run setup commands |
| `womm init [dir]` | Generate a starter `.womm` file, optionally inside a new directory |
| `womm verify` | Read-only check — reports what's missing without installing or running anything |
| `womm -c <path>` | Use a specific config file instead of `.womm` |

## Full Example

See [`examples/setup.womm`](examples/setup.womm) for a complete configuration reference.

```sh
# Run against a specific config
womm -c examples/setup.womm

# Verify a remote project's requirements
womm verify -c path/to/.womm
```

## How It Works

1. **Resolve config** — Searches for `.womm` in the current directory, or uses the path from `-c`.
2. **Platform check** — Validates that the current OS is in the allowed list.
3. **Dependency check** — Runs `<tool> --version` for each entry and parses the output against the required version.
4. **Environment check** — Looks up each variable in `os.Environ()`.
5. **File check** — Verifies each path exists on disk.
6. **Install** (default only) — For missing or outdated dependencies, detects available package managers and prompts the user to install.
7. **Commands** (default only) — Prompts before executing each command group.
