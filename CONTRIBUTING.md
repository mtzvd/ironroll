# Contributing

Contributions are welcome.

## Development

### Prerequisites

- Go 1.21+
- [Task](https://taskfile.dev/) (optional)

### Setup

```bash
git clone https://github.com/mtzvd/ironroll.git
cd ironroll
go mod download
```

### Running Tests

```bash
go test ./...
```

Or with Task:

```bash
task test
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Keep the core/roll package free of external dependencies

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## Issues

Feel free to open issues for bugs, feature requests, or questions.
