# 🛒 E-commerce Microservices & Event-Driven Architecture

Este projeto é um ecossistema completo de e-commerce baseado em **microserviços**, demonstrando uma arquitetura robusta, escalável e desacoplada. O sistema utiliza **Go** para o back-end, **React** para o front-end e **RabbitMQ** para a comunicação assíncrona baseada em eventos, além de implementar o padrão de **Webhooks** para integrações externas.

<img width="778" height="710" alt="Captura de Tela 2026-03-09 às 17 31 08" src="https://github.com/user-attachments/assets/4acffc4a-4e91-4adf-8bf3-9165932c8217" />

---

## 🏗️ Arquitetura do Sistema

O projeto é dividido em serviços especializados, cada um com sua própria responsabilidade e banco de dados (**Database per Service**):

### 🔙 Back-end (Go)
1.  **Ecommerce API (`ecommerce-api`):** Gerencia o catálogo de produtos e o ciclo de vida dos pedidos.
2.  **Admin Ecommerce API (`api-admin-ecommerce`):** Atua como o "back-office", processando solicitações de pagamento e orquestrando a comunicação entre o e-commerce e o gateway.
3.  **Payment Gateway API (`gateway-api`):** Simula um gateway de pagamento de mercado (como Stripe ou Pagar.me), processando transações e notificando o sistema via Webhooks.

### 🎨 Front-end (React + TailwindCSS)
1.  **Storefront:** Interface para o cliente final realizar compras.
2.  **Admin Panel:** Interface para gestão de pagamentos e pedidos.
3.  **Gateway Panel:** Interface administrativa do gateway para simular aprovações/rejeições de pagamentos.

---

## 🚀 Tecnologias Utilizadas

-   **Linguagens:** Go (Golang), JavaScript (React).
-   **Mensageria:** RabbitMQ (AMQP 0.9.1).
-   **Banco de Dados:** MySQL 8.0.
-   **Containerização:** Docker & Docker Compose.
-   **Arquitetura:** Clean Architecture, Event-Driven Design, Microservices.
-   **Comunicação:** REST (HTTP), Webhooks, Mensageria Assíncrona.

---

## 🔄 Fluxo de um Pedido (Event-Driven)

O sistema demonstra o poder do processamento assíncrono para garantir resiliência:

1.  **Criação do Pedido:** O cliente finaliza a compra no React Storefront. A `ecommerce-api` salva o pedido como `pending`.
2.  **Solicitação de Pagamento:** A `ecommerce-api` publica um evento `payment.requested` no **RabbitMQ**.
3.  **Orquestração:** A `api-admin-ecommerce` consome esse evento, registra a intenção de pagamento e envia os dados para a `gateway-api`.
4.  **Processamento do Gateway:** O gateway processa a transação. Assim que o status muda, ele dispara um **Webhook** de volta para a `api-admin-ecommerce`.
5.  **Notificação de Resultado:** A `api-admin-ecommerce` recebe o webhook, atualiza seu registro local e publica um evento `payment.processed` no **RabbitMQ**.
6.  **Finalização:** A `ecommerce-api` consome o evento final e atualiza o status do pedido para `paid` ou `rejected`, refletindo instantaneamente para o cliente.

---

## 🛠️ Como Executar

### Pré-requisitos
-   Docker e Docker Compose instalados.
-   Go 1.22+ (opcional, para desenvolvimento local).
-   Node.js & npm (opcional, para os front-ends).

### Passo a Passo
1.  **Subir a Infraestrutura:**
    ```bash
    docker-compose up -d --build
    ```
    Isso iniciará:
    - RabbitMQ (Broker)
    - 3 Instâncias de MySQL
    - 3 APIs em Go
    - PHPMyAdmin (Gerenciamento de DB)

2.  **Acessar os Serviços:**
    - **Ecommerce API:** `http://localhost:3000`
    - **Admin API:** `http://localhost:9000`
    - **Gateway API:** `http://localhost:4000`
    - **RabbitMQ Management:** `http://localhost:15672` (guest/guest)
    - **PHPMyAdmin:** `http://localhost:8081`

3.  **Executar os Front-ends:**
    Entre em cada pasta de frontend (`frontend-ecommerce`, `admin-ecommerce`, `gateway-admin-panel`) e execute:
    ```bash
    npm install
    npm run dev
    ```

---

## 💡 Destaques de Implementação

-   **Resiliência:** Se um serviço cair, as mensagens ficam retidas no RabbitMQ e são processadas assim que o serviço retornar.
-   **Desacoplamento:** A `ecommerce-api` não sabe da existência da `gateway-api`. Elas se comunicam estritamente por contratos de mensagens.
-   **Segurança e Logs:** Implementação de logs estruturados e tratamento de erros robusto em Go.
-   **Clean Architecture:** Separação clara entre Entidades, Casos de Uso e Adaptadores de Infraestrutura, facilitando testes e manutenção.

---

Desenvolvido para estudo de sistemas distribuídos e mensageria. 🚀
