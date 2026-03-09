# Payment Gateway Example

This project is a simple payment gateway API built with Go, demonstrating a clean architecture approach with layered components. It supports basic CRUD operations for payments, including creation, retrieval (single and paginated list), update, and deletion.

## Technical Overview

The application follows a clean architecture pattern, separating concerns into distinct layers:

*   **Domain Layer (`internal/domain`)**: Contains the core business logic, entities (`entity/payment.go`), and repository interfaces (`repository/payment_repository.go`). This layer is independent of external frameworks or databases.
*   **Use Case Layer (`internal/usecase`)**: Implements the application's specific business rules and orchestrates interactions between the domain and interface layers. Each use case represents a single operation (e.g., `create_payment.go`, `get_payment.go`, `get_all_payments.go`, `update_payment.go`, `delete_payment.go`).
*   **Infrastructure Layer (`internal/infrastructure`)**: Provides the implementations for external concerns such as database persistence (`database/mysql/payment_repository.go`) and configuration (`config/config.go`).
*   **Interface Layer (`internal/interface`)**: Handles external communication, including HTTP requests and responses (`http/handler/payment_handler.go`, `http/router.go`) and Data Transfer Objects (DTOs) (`dto/payment_dto.go`).

### Technologies Used

*   **Go (Golang)**: The primary programming language.
*   **Chi Router**: A lightweight, idiomatic, and composable router for building HTTP services in Go.
*   **MySQL**: The relational database used for storing payment information.
*   **Docker**: Used for containerization of the application and database.
*   **Docker Compose**: For defining and running multi-container Docker applications.

## Key Components

*   **`payment.go` (Entity)**: Defines the `Payment` struct and its behavior.
*   **`payment_repository.go` (Repository Interface)**: Declares the contract for data persistence operations (Save, FindByID, FindAll, Delete).
*   **`mysql/payment_repository.go` (MySQL Repository Implementation)**: Provides the concrete implementation of the `PaymentRepository` interface using MySQL.
*   **Use Cases**:
    *   `CreatePayment`: Handles the creation of new payments.
    *   `GetPayment`: Retrieves a single payment by ID.
    *   `GetAllPayments`: Retrieves a paginated list of payments.
    *   `UpdatePayment`: Updates the status of an existing payment.
    *   `DeletePayment`: Deletes a payment by ID.
*   **`payment_handler.go`**: Contains HTTP handlers for each payment-related operation, parsing requests, calling the appropriate use cases, and formatting responses.
*   **`router.go`**: Sets up the HTTP routes and associates them with the handlers.
*   **`main.go`**: The application entry point, responsible for initializing the database connection, repositories, use cases, and HTTP router.

## Setup and Running

### Prerequisites

*   Go (version 1.18 or higher)
*   Docker
*   Docker Compose

### Steps

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/gateway-payments-example.git
    cd gateway-payments-example
    ```

2.  **Start Docker Compose services:**
    This will build the Go application image, start the MySQL database, and the Nginx reverse proxy.
    ```bash
    docker-compose up --build -d
    ```

3.  **Initialize the database:**
    Execute the SQL script to create the `payments` table. You can do this by connecting to the MySQL container or by using a tool like `mysql` CLI.
    ```bash
    docker exec -i gateway-payments-mysql mysql -uroot -pmysecretpassword payments < create_table.sql
    ```
    (Replace `gateway-payments-mysql` with the actual name of your MySQL service container if it's different).

4.  **Access the application:**
    The application will be accessible via the Nginx reverse proxy.
    *   **Base URL**: `http://localhost:80` (or `http://localhost:8080` if accessing the Go app directly)

## API Endpoints

All endpoints are prefixed with `/payments`.

*   **`POST /payments`**: Create a new payment.
    *   Request Body: `{"method": "CreditCard", "amount": 100.00}`
    *   Response: `201 Created` with the created payment details.

*   **`GET /payments/{id}`**: Retrieve a single payment by ID.
    *   Response: `200 OK` with the payment details, or `404 Not Found`.

*   **`GET /payments?page={page}&limit={limit}`**: Retrieve a paginated list of payments.
    *   Query Parameters:
        *   `page` (optional, default 1): The page number.
        *   `limit` (optional, default 10): The number of items per page.
    *   Response: `200 OK` with an array of payment details.

*   **`PUT /payments/{id}`**: Update the status of a payment.
    *   Request Body: `{"status": "approved"}`
    *   Response: `200 OK` or `404 Not Found`.

*   **`DELETE /payments/{id}`**: Delete a payment by ID.
    *   Response: `204 No Content` or `404 Not Found`.
