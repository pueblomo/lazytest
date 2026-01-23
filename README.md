# lazytest

A simple terminal UI for running tests, inspired by [lazydocker](https://github.com/jesseduffield/lazydocker) and [lazygit](https://github.com/jesseduffield/lazygit).

## Why lazytest?

Because running tests shouldn't require remembering complex commands or switching between terminal windows. **lazytest** gives you a keyboard-driven interface to:

- 🔍 **Discover tests** automatically in your project
- ▶️ **Run tests** with a single keystroke
- 👀 **View output** in real-time
- 🔄 **Re-run tests** instantly as you fix issues
- 📊 **Track status**

## Features

- **Auto-detection**: Automatically detects your test framework (currently supports Vitest)
- **Interactive UI**: Navigate tests with vim-style keybindings
- **Live output**: See test results as they happen
- **Fast filtering**: Quickly find tests with fuzzy search
- **Watch mode**: Auto-detects file changes

## Installation

### From source

```bash
go install github.com/pueblomo/lazytest/cmd/lazytest@latest
```

### Manual

```bash
git clone https://github.com/pueblomo/lazytest.git
cd lazytest
go build -o lazytest ./cmd/lazytest
```

## Usage

Navigate to your project directory and run:

```bash
lazytest
```

### Keybindings

#### Test List Panel
- `↑/k` - Move up
- `↓/j` - Move down
- `enter` - Run selected test
- `/` - Filter tests
- `tab` - Switch focus
- `w` - Toggle watch mode for selected test
- `q` - Quit

#### Output Panel
- `tab` - Switch focus
- `↑/k` - Scroll up
- `↓/j` - Scroll down
- `q` - Quit

#### Logs Panel
- `tab` - Switch focus
- `↑/k` - Scroll up
- `↓/j` - Scroll down
- `q` - Quit

## Supported Test Frameworks

- [x] Vitest
- [ ] Jest
- [ ] Pytest
- [ ] Go test
- [ ] RSpec
- [ ] PHPUnit
- [ ] Your framework here? [Submit a PR!](CONTRIBUTING.md)

## Requirements

- Go 1.25+ (for building from source)
- A supported test framework installed in your project

## Contributing

Contributions are welcome! The goal is to support as many test frameworks as possible with a pluggable driver architecture.

## License

MIT

## Inspiration

- [lazygit](https://github.com/jesseduffield/lazygit) - A simple terminal UI for git commands
- [lazydocker](https://github.com/jesseduffield/lazydocker) - A simple terminal UI for docker commands
