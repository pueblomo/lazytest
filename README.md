# lazytest

A simple terminal UI for running tests, inspired by [lazydocker](https://github.com/jesseduffield/lazydocker) and [lazygit](https://github.com/jesseduffield/lazygit).

<p align="center">
  <a href="https://github.com/pueblomo/lazytest/actions/workflows/test.yml/badge.svg"><img src="https://github.com/pueblomo/lazytest/actions/workflows/test.yml/badge.svg" alt="Test Badge"></a>
  <a href="https://codecov.io/gh/pueblomo/lazytest/branch/main/graph/badge.svg"><img src="https://codecov.io/gh/pueblomo/lazytest/branch/main/graph/badge.svg" alt="Coverage Badge"></a>
</p>

## Why lazytest?

Running tests shouldn't require memorizing commands, switching windows, or losing focus. **lazytest** brings your test suite into a single, keyboard-driven interface so you can:

- Discover tests automatically across your project
- Run, filter, and re-run tests without leaving your terminal
- Watch live output with clear pass/fail indicators
- Stay in flow with vim-style navigation and minimal context switching

## Features

- **Multi-framework support** out of the box (Go, Vitest, Maven, Gradle)
- **Auto-detection** of your project's test runner
- **Live results** with color-coded output
- **Interactive filtering** (fuzzy search) to find tests fast
- **Watch mode** that re-runs affected tests on file changes
- **Keyboard-centric** UI with helpful shortcuts (press `?` in-app)
- **Cross-platform**: Linux, macOS, and Windows

## Supported Test Frameworks

- ✅ **Go** (`go test`)
- ✅ **Vitest** (for JavaScript/TypeScript projects)
- ✅ **Maven** (`mvn test`)
- ✅ **Gradle** (`./gradlew test`)

Planning to add more drivers? See the [Contributing](#contributing) section.

## Installation

### From source (requires Go 1.25+)

```bash
go install github.com/pueblomo/lazytest/cmd/lazytest@latest
```

### Pre-built binaries

Download the latest release for your platform from the [GitHub Releases page](https://github.com/pueblomo/lazytest/releases).

### Manual build

```bash
git clone https://github.com/pueblomo/lazytest.git
cd lazytest
go build -o lazytest ./cmd/lazytest
```

### Alternative package managers (community)

Homebrew, Scoop, and other package manager taps may become available—check the repository for up-to-date instructions.

## Usage

Navigate to your project directory and start the UI:

```bash
lazytest
```

Use arrow keys or `hjkl` to navigate, `Enter` to run a test, `r` to re-run, and `/` to filter. Press `?` for the full list of shortcuts.

## Platform Support

lazytest is a single Go binary that runs on:

- Linux (amd64, arm64)
- macOS (Intel, Apple Silicon)
- Windows (amd64)

## Contributing

Contributions are welcome! The project is built with a **pluggable driver architecture**—each test framework is implemented as a driver under `internal/drivers`.

To contribute:

1. Fork the repository and create a feature branch.
2. Follow the existing driver pattern if adding a new test framework.
3. Include tests for your changes (`go test ./...`).
4. Ensure `go fmt` and `go vet` pass.
5. Open a pull request with a clear description.

Please see the in-app documentation (`?`) or the code comments for deeper implementation details. For larger changes, consider opening an issue first to discuss design.

## License

MIT

## Inspiration

- [lazygit](https://github.com/jesseduffield/lazygit) – A simple terminal UI for git commands
- [lazydocker](https://github.com/jesseduffield/lazydocker) – A simple terminal UI for docker commands

---

> **Note:** Screenshots and GIFs will be added to this section once available to showcase the UI.
