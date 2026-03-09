# E-commerce API Example

This project is an example of an e-commerce API built with Go, using a clean architecture approach. It provides functionalities for managing products (items) and customer orders.

## Technologies Used

*   **Go**: The primary programming language.
*   **MySQL**: Relational database for data persistence.
*   **Docker**: For containerization of the application and database.
*   **Docker Compose**: For defining and running multi-container Docker applications.
*   **Nginx**: As a reverse proxy.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Before you begin, ensure you have the following installed on your system:

*   [Docker](https://docs.docker.com/get-docker/)
*   [Docker Compose](https://docs.docker.com/compose/install/) (usually comes with Docker Desktop)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/api-ecommerce-example.git
    cd api-ecommerce-example
    ```

2.  **Set up the database:**
    The project uses MySQL. The `docker-compose.yml` file will set up a MySQL container. After the database container is running, you'll need to create the necessary tables.

    First, bring up the database service:
    ```bash
    docker-compose up -d db
    ```

    Create a `.env` file at the root of the project with the following content (replace with your desired values):

    ```
    API_PORT=8080
    DB_PORT=3306
    MYSQL_ROOT_PASSWORD=root_password
    MYSQL_DATABASE=ecommerce_db
    MYSQL_USER=user
    MYSQL_PASSWORD=password
    ```

    Once the MySQL container is up and running, and your `.env` file is configured, you can execute the `create_table.sql` script.

    ```bash
    docker exec -i ecommerce-api-db mysql -u root -p"$MYSQL_ROOT_PASSWORD" < create_table.sql
    ```
    (Ensure you replace `$MYSQL_ROOT_PASSWORD` with the actual root password you set in your `.env` file or provide it directly after `-p` without a space).

### Running the Application

1.  **Build and run all services using Docker Compose:**

    ```bash
    docker-compose up --build
    ```

    This command will build the Go application image, start the MySQL database, and the Nginx reverse proxy.

2.  **Access the API:**
    The API should be accessible via Nginx, typically at `http://localhost:80`.

## API Endpoints

The API provides endpoints for managing `items` and `orders`.

### Items

*   `POST /items`: Create a new item.
*   `GET /items`: Get all items.
*   `GET /items/{id}`: Get an item by ID.
*   `PUT /items/{id}`: Update an item by ID.
*   `DELETE /items/{id}`: Delete an item by ID.

### Orders

*   `POST /orders`: Create a new order.
*   `GET /orders`: Get all orders.
*   `GET /orders/{id}`: Get an order by ID.

## Project Structure

*   `cmd/api`: Contains the main application entry point.
*   `internal/`: Core application logic, divided into:
    *   `domain/`: Business entities and interfaces (repositories).
    *   `infrastructure/`: Implementations of interfaces, database connections, and configurations.
    *   `interface/`: HTTP handlers and DTOs for external communication.
    *   `usecase/`: Application-specific business rules and orchestrators.
*   `nginx/`: Nginx configuration files.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
