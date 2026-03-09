# 🛒 Go & React E-Commerce API

Um projeto de e-commerce completo utilizando **Clean Architecture** em Go e um frontend moderno em **React + Tailwind CSS**.

## 🚀 Tecnologias
- **Backend:** Go (Golang) com roteamento nativo.
- **Frontend:** React + Vite + Tailwind CSS v4.
- **Mensageria:** RabbitMQ para processamento de pedidos.
- **Banco de Dados:** MySQL.
- **Container:** Docker & Nginx.

## 🏗️ Arquitetura do Sistema
O sistema funciona de forma assíncrona:
1. O Frontend envia um pedido para a API Go.
2. A API salva o pedido no MySQL com status `pending`.
3. Um evento é disparado para o **RabbitMQ**.
4. Um worker consome a fila e atualiza o status para `paid` ou `approved`.



## 🛠️ Como Rodar
1. Suba os serviços (MySQL/RabbitMQ): `docker-compose up -d`
2. Rode a API: `go run cmd/api/main.go`
3. Rode o Front: 
   - `cd frontend-ecommerce`
   - `npm install && npm run dev`