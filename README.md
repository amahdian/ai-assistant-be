# AI Assistant Backend

This project is the backend for a powerful and flexible AI assistant that connects to various Large Language Models (LLMs). It's built with Go and the Gin framework, providing a robust and scalable foundation for building chat-based AI applications. This backend is designed to work seamlessly with the [AI Assistant Flutter app](https://github.com/amahdian/ai_assistant_flutter).

## 🚀 Features

*   **📱 Cross-Platform Support**: Works with the AI Assistant mobile app for both Android and iOS.
*   **Multiple LLM Support**: Easily integrate with different LLMs to power your chat assistant.
*   **Secure Authentication**: Built-in JWT-based authentication to protect your users' data.
*   **Scalable Architecture**: Designed for performance and scalability, ready for production use.
*   **Comprehensive API Documentation**: Auto-generated Swagger documentation for easy API exploration.
*   **Dockerized Environment**: Comes with a Docker Compose setup for a smooth development experience.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

*   **Go 1.24.2+**: [Download Go](https://golang.org/dl/)
*   **PostgreSQL 10.3+**: [Download PostgreSQL](https://www.postgresql.org/download/)
*   **Docker & Docker Compose** (recommended): [Download Docker](https://www.docker.com/products/docker-desktop)
*   **Make**: For using the provided Makefile commands.
*   **golang-migrate**: For database migrations.

### Installing golang-migrate

You can install `golang-migrate` using one of the following methods:

*   **Go Install (Recommended)**:
    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
*   **Homebrew (macOS)**:
    ```bash
    brew install golang-migrate
    ```

## 🛠️ Setup

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/amahdian/ai-assistant-be.git
    cd ai-assistant-be
    ```

2.  **Environment Configuration**:
    Create an environment file by copying the example:
    ```bash
    cp .env.example .env
    ```
    Update the `.env` file with your configuration, especially the database connection string and JWT secret.

3.  **Database Setup**:
    *   **Using Docker (Recommended)**:
        ```bash
        docker-compose up -d postgres
        make create-db
        make migrate-up
        ```
    *   **Using Local PostgreSQL**:
        ```bash
        # Create the database manually
        createdb app_db

        # Run migrations
        make migrate-up
        ```

4.  **Install Dependencies**:
    ```bash
    make vendor
    ```

## 🚀 Running the Application

*   **Development Mode**:
    ```bash
    make dev
    ```
*   **Production Mode**:
    ```bash
    make build
    ./build/app-bin
    ```

## 📖 API Documentation

Once the application is running, you can access the Swagger UI for API documentation at:

[http://localhost:8090/swagger/index.html](http://localhost:8090/swagger/index.html)

## 🏗️ Project Structure

The project follows a standard Go project layout:

```
ai-assistant-be/
├── assets/         # Static assets and migrations
├── docs/           # Swagger documentation
├── domain/         # Domain models and contracts
├── global/         # Global configurations and utilities
├── pkg/            # Reusable packages
├── server/         # HTTP server components
├── storage/        # Data storage layer
├── svc/            # Business logic services
├── main.go         # Application entry point
└── Makefile        # Build and development commands
```

## 🤝 Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature`).
3.  Commit your changes (`git commit -m 'Add some feature'`).
4.  Push to the branch (`git push origin feature/your-feature`).
5.  Open a pull request.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.