# GoFiber API with PostgreSQL and Redis

A high-performance API built with Go using the GoFiber framework, integrating PostgreSQL for data storage and Redis for caching. This project provides a clean structure and efficient data querying, suitable for scalable applications.

## Features

- **RESTful API**: Simple and intuitive API endpoints for user data retrieval.
- **Caching**: Uses Redis for caching database queries, improving performance.
- **Database Integration**: PostgreSQL connection with GORM for ORM support.
- **Configuration Management**: JSON-based configuration with `viper` for flexible settings.
- **Logging**: Comprehensive logging middleware for tracking requests and errors.
- **Error Handling**: Unified error handling middleware for consistent responses.

## Getting Started

### Prerequisites

- Go (1.16+)
- PostgreSQL
- Redis

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/iyuangang/go-fiber-api.git
   cd go-fiber-api
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up your PostgreSQL and Redis databases and update the configuration in `config/config.json`.

### Configuration

Edit the `config/config.json` file:

```json
{
  "postgres": {
    "url": "postgres://user:password@localhost:5432/mydb",
    "max_idle_conns": 10,
    "max_open_conns": 100,
    "conn_max_lifetime": 300
  },
  "redis": {
    "addr": "localhost:6379",
    "pass": "",
    "db": 0,
    "cache_expiration_minutes": 10
  },
  "server": {
    "port": 3000,
    "read_timeout": 5
  }
}
```

### Running the Application

Run the application:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:3000`.

### API Endpoints

- **GET /user/:id**: Retrieve user information by ID. Caches results in Redis.

### Example Request

```bash
curl -X GET http://localhost:3000/user/1
```

### Logging

The application logs request details and errors to the console. Adjust logging levels in the application for production environments.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any improvements or bug reports.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [GoFiber](https://gofiber.io) - Fast web framework for Go.
- [GORM](https://gorm.io) - ORM library for Go.
- [Redis](https://redis.io) - In-memory data structure store.
```
