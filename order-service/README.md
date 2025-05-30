# Order Service

This microservice handles order creation, delivery address management, and delivery time selection for the online grocery store.

## Features

- Order creation and management
- Delivery address CRUD operations
- Delivery time selection
- Integration with cart service
- Real-time order status updates via NATS

## Tech Stack

- Go 1.21+
- gRPC
- MongoDB
- Redis
- NATS
- JWT Authentication
- Zap Logger
- Wire (Dependency Injection)
- Docker

## Project Structure

```
order-service/
├── cmd/                    # Application entry points
├── internal/              
│   ├── domain/            # Enterprise business rules
│   ├── usecase/           # Application business rules
│   ├── repository/        # Data access implementations
│   ├── delivery/          # Delivery mechanisms (gRPC, HTTP)
│   └── infrastructure/    # External services, DB, cache
├── pkg/                   # Public packages
├── proto/                 # Protocol buffer definitions
├── migrations/            # Database migrations
├── config/               # Configuration files
└── test/                 # Integration tests
```

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
```bash
cp .env.example .env
```

3. Start required services:
```bash
docker-compose up -d
```

4. Run database migrations:
```bash
make migrate-up
```

5. Start the service:
```bash
make run
```

## Development

### Running Tests
```bash
make test
```

### Generate Proto Files
```bash
make proto
```

### Database Migrations
```bash
# Create new migration
make migrate-create name=migration_name

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## API Documentation

See `proto/order.proto` for the complete API specification.

## Monitoring

The service exposes metrics at `/metrics` for Prometheus scraping.

## License

MIT 