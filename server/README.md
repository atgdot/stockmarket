# Stock Market Application

A robust and scalable stock market application built with Go, featuring real-time data processing, secure authentication, and cloud-native architecture.

## 🚀 Technologies & Skills Demonstrated

### Backend Development
- **Go (Golang)** - Core application development
- **Echo Framework** - High-performance HTTP server and routing
- **JWT Authentication** - Secure user authentication and authorization
- **WebSocket** - Real-time data streaming and updates

### Cloud & Infrastructure
- **AWS Services**
  - DynamoDB - NoSQL database for flexible data storage
  - SNS (Simple Notification Service) - Event-driven messaging
  - AWS SDK v2 - Modern cloud integration
- **Redis** - In-memory data store for caching and real-time features

### Architecture & Design
- Clean Architecture principles
- Microservices-ready design
- Environment-based configuration
- Modular code organization

### Development Practices
- Dependency management with Go modules
- Environment variable management
- API documentation
- Postman collection for API testing

## 🛠️ Prerequisites

- Go 1.24 or higher
- AWS Account with appropriate credentials
- Redis server
- Docker (optional, for containerization)

## 🔧 Setup Instructions

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd stockmarket/server
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   Create a `.env` file in the root directory with the following variables:
   ```
   AWS_ACCESS_KEY_ID=your_access_key
   AWS_SECRET_ACCESS_KEY=your_secret_key
   AWS_REGION=your_region
   REDIS_URL=your_redis_url
   JWT_SECRET=your_jwt_secret
   ```

4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

   For development with hot-reload:
   ```bash
   air
   ```

## 📚 API Documentation

API documentation is available in the `postman` directory. Import the Postman collection to test the endpoints.

## 🏗️ Project Structure

```
server/
├── api/            # API handlers and routes
├── cmd/            # Application entry points
├── configs/        # Configuration files
├── deployment/     # Deployment configurations
├── internal/       # Internal packages
│   ├── database/   # Database interactions
│   └── ...         # Other internal packages
├── pkg/            # Public packages
└── postman/        # API documentation
```

## 🔐 Security Features

- JWT-based authentication
- Secure password hashing
- Environment variable management
- AWS IAM integration

## 🚀 Performance Optimizations

- Redis caching
- WebSocket for real-time updates
- Efficient database queries
- Connection pooling

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

[Your Name]

## 🙏 Acknowledgments

- AWS SDK team for the excellent Go SDK
- Echo framework contributors
- Redis team for the Go client
