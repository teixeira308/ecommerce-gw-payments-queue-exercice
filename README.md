# 🛒 E-commerce Microservices & Event-Driven Architecture

This project is a complete e-commerce ecosystem based on **microservices**, showcasing a robust, scalable, and decoupled architecture. The system uses **Go** for the back end, **React** for the front end, and **RabbitMQ** for asynchronous event-based communication, while also implementing the **Webhook** pattern for external integrations.


<img width="778" height="702" alt="Captura de Tela 2026-03-11 às 18 58 08" src="https://github.com/user-attachments/assets/ab73121b-2a2a-4ff0-b92b-ee01ccc7b882" />


---

## 🏗️ System Architecture

The project is divided into specialized services, each with its own responsibility and database (**Database per Service**):

### 🔙 Back-end (Go)
1.  **Ecommerce API (`ecommerce-api`):** Manages the product catalog and the order lifecycle.
2.  **Admin Ecommerce API (`api-admin-ecommerce`):** Acts as the back office, processing payment requests and orchestrating communication between the e-commerce system and the gateway.
3.  **Payment Gateway API (`gateway-api`):** Simulates a market payment gateway (such as Stripe or Pagar.me), processing transactions and notifying the system via Webhooks.

### 🎨 Front-end (React + TailwindCSS)
1.  **Storefront:** Interface for end customers to place orders.
2.  **Admin Panel:** Interface for managing payments and orders.
3.  **Gateway Panel:** Gateway admin interface to simulate payment approvals and rejections.

---

## 🚀 Technologies Used

-   **Languages:** Go (Golang), JavaScript (React).
-   **Messaging:** RabbitMQ (AMQP 0.9.1).
-   **Database:** MySQL 8.0.
-   **Containerization:** Docker & Docker Compose.
-   **Architecture:** Clean Architecture, Event-Driven Design, Microservices.
-   **Communication:** REST (HTTP), Webhooks, Asynchronous Messaging.

---

## 🔄 Order Flow (Event-Driven)

The system demonstrates the power of asynchronous processing to ensure resilience:

1.  **Order Creation:** The customer completes the purchase in the React Storefront. The `ecommerce-api` saves the order as `pending`.
2.  **Payment Request:** The `ecommerce-api` publishes a `payment.requested` event to **RabbitMQ**.
3.  **Orchestration:** The `api-admin-ecommerce` consumes this event, records the payment intent, and sends the data to the `gateway-api`.
4.  **Gateway Processing:** The gateway processes the transaction. As soon as the status changes, it sends a **Webhook** back to the `api-admin-ecommerce`.
5.  **Result Notification:** The `api-admin-ecommerce` receives the webhook, updates its local record, and publishes a `payment.processed` event to **RabbitMQ**.
6.  **Completion:** The `ecommerce-api` consumes the final event and updates the order status to `paid` or `rejected`, reflecting the result instantly to the customer.

---

## 🛠️ How to Run

### Prerequisites
-   Docker and Docker Compose installed.
-   Go 1.22+ (optional, for local development).
-   Node.js & npm (optional, for the front ends).

### Step by Step
1.  **Start the Infrastructure:**
    ```bash
    docker-compose up -d --build
    ```
    This will start:
    - RabbitMQ (Broker)
    - 3 MySQL instances
    - 3 Go APIs
    - PHPMyAdmin (DB management)

2.  **Access the Services:**
    - **Ecommerce API:** `http://localhost:3000`
    - **Admin API:** `http://localhost:9000`
    - **Gateway API:** `http://localhost:4000`
    - **RabbitMQ Management:** `http://localhost:15672` (guest/guest)
    - **PHPMyAdmin:** `http://localhost:8081`

3.  **Run the Front Ends:**
    Go into each frontend folder (`frontend-ecommerce`, `admin-ecommerce`, `gateway-admin-panel`) and run:
    ```bash
    npm install
    npm run dev
    ```

---

## 💡 Implementation Highlights

-   **Resilience:** If a service goes down, messages remain queued in RabbitMQ and are processed as soon as the service comes back.
-   **Decoupling:** The `ecommerce-api` is unaware of the `gateway-api`. They communicate strictly through message contracts.
-   **Security and Logging:** Structured logging and robust error handling implemented in Go.
-   **Clean Architecture:** Clear separation between Entities, Use Cases, and Infrastructure Adapters, making testing and maintenance easier.

---

Built for studying distributed systems and messaging. 🚀
