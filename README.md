# Yurt Mart - Microservices E-commerce Platform

## Purpose
Yurt Mart is the final project for our course, designed as a modern, scalable e-commerce platform using a microservices architecture. The project demonstrates advanced concepts in distributed systems, service communication, and cloud-native deployment using Docker.

## Architecture Overview
This project is composed of several independent microservices, each responsible for a specific domain of the e-commerce platform. The services communicate via gRPC, HTTP, and NATS messaging, and are orchestrated using Docker Compose for easy local development and deployment.

### Main Services
- **API Gateway**: Entry point for all client requests, routing traffic to appropriate backend services.
- **User Service**: Manages user registration, authentication, and profiles.
- **Product Service**: Handles product catalog management and queries.
- **Order Service**: Processes customer orders, manages order lifecycle.
- **Payment Service**: Handles payment processing, including integration with blockchain and credit card processors.
- **Shopping Cart Service**: Manages user shopping carts and cart operations.
- **Order History Service**: Stores and retrieves historical order data.
- **Review Service**: Manages product reviews and ratings.

## Technologies Used
- **Go (Golang)** for all backend services
- **MongoDB** for data persistence
- **NATS** for event-driven communication
- **gRPC & HTTP** for service APIs
- **Docker & Docker Compose** for containerization and orchestration

## Getting Started (Docker)

### Prerequisites
- [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/) installed on your machine

### Running the Project
1. Clone this repository:
   ```sh
   git clone <your-repo-url>
   cd Yurt_Mart-master (2)
   ```
2. Build and start all services using Docker Compose:
   ```sh
   docker-compose up --build
   ```
3. Access the API Gateway at `http://localhost:<gateway-port>` (replace `<gateway-port>` with the actual port in your docker-compose.yml)

### Stopping the Project
```sh
docker-compose down
```

## Service Endpoints
Each service exposes its own gRPC and/or HTTP endpoints. See the `proto/` directories and service documentation for details.



## License
- *Specify your license here (e.g., MIT, Apache 2.0, etc.)*

---
*This project is developed as a final project for educational purposes.* 
