# Payment Service

This microservice handles payment processing for the online grocery store, supporting both credit card and MetaMask (crypto) payments.

## Features

- Credit/Debit card payment processing
- MetaMask (crypto) payment integration
- Payment status tracking
- Payment recovery for incomplete transactions
- Real-time payment status updates via NATS

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
- Web3 Integration (for MetaMask)

## Project Structure

```
payment-service/
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

## Smart Contract Integration

The payment service integrates with Ethereum smart contracts for crypto payments. See `contracts/` directory for smart contract implementations.

## API Documentation

See `proto/payment.proto` for the complete API specification.

## Monitoring

The service exposes metrics at `/metrics` for Prometheus scraping.

## License

MIT 