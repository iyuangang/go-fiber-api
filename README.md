# Go Fiber API

A high-performance RESTful API built with Go and the Fiber web framework.

## Features

- Fast and efficient API endpoints using Fiber
- PostgreSQL database integration with GORM
- Redis caching for improved performance
- Structured logging with Zap
- Containerization with Docker
- Multi-platform build support (Linux, macOS, Windows)
- Continuous Integration and Deployment with GitHub Actions

## Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Redis
- Docker (optional)

## Getting Started

### Clone the repository

```bash
git clone https://github.com/iyuangang/go-fiber-api.git
cd go-fiber-api
```

### Set up the configuration

Copy the example configuration file and adjust it to your needs:

```bash
cp config/config.example.yaml config/config.yaml
```

Edit `config/config.yaml` with your database and Redis credentials.

### Run the application

```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:3000`.

## API Endpoints

- `GET /api/users`: Get all users
- `GET /api/user/:id`: Get a user by ID
- `POST /api/user`: Create a new user
- `PUT /api/user/:id`: Update a user
- `DELETE /api/user/:id`: Delete a user

For detailed API documentation, please refer to the [API Documentation](docs/api.md).

### Example Request

```bash
curl -X GET http://localhost:3000/user/1
```

## Development

### Running tests

```bash
go test ./...
```

### Building for different platforms

The project includes GitHub Actions workflows for building the application for Linux, macOS, and Windows. You can also build manually:

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o go-fiber-api-linux-amd64 ./cmd/server

# For macOS
GOOS=darwin GOARCH=amd64 go build -o go-fiber-api-darwin-amd64 ./cmd/server

# For Windows
GOOS=windows GOARCH=amd64 go build -o go-fiber-api-windows-amd64.exe ./cmd/server
```

## Deployment

### Using Docker

Build the Docker image:

```bash
docker build -t go-fiber-api .
```

Run the container:

```bash
docker run -p 3000:3000 go-fiber-api
```

### CI/CD

The project uses GitHub Actions for Continuous Integration and Deployment. The workflows are defined in `.github/workflows/`:

- Runs tests and builds the application for the `dev` branch.
- Builds and publishes Docker images for releases on the `main` branch.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [GoFiber](https://gofiber.io) - Fast web framework for Go.
- [GORM](https://gorm.io) - ORM library for Go.
- [Redis](https://redis.io) - In-memory data structure store.
