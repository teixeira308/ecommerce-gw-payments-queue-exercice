CREATE DATABASE IF NOT EXISTS gateway_db;
USE gateway_db;

CREATE TABLE IF NOT EXISTS payments (
    -- ID v4 do UUID (36 caracteres)
    id CHAR(36) NOT NULL,
    
    -- Valor financeiro com 2 casas decimais
    amount DECIMAL(10, 2) NOT NULL,
    
    -- Método de pagamento (PIX, CREDIT_CARD, etc)
    method VARCHAR(20) NOT NULL,

    -- ID da Order associada a esse pagamento
    order_id VARCHAR(36) NULL,
    
    -- Status: Aceita NULL 
    status VARCHAR(20) NULL,
    
    -- Data de criação com precisão de microsegundos
    created_at DATETIME(6) NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_status (status),
    INDEX idx_order_id (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;