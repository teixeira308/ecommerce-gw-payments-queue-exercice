💳 Gateway Admin Panel
Este é o painel administrativo do Gateway de Pagamentos. Desenvolvido em React, ele permite que administradores monitorem transações em tempo real, aprovem ou rejeitem pagamentos pendentes e acompanhem métricas de faturamento.

🚀 Funcionalidades
Dashboard Financeiro: Visualização do faturamento total aprovado e volume de transações pendentes.

Gestão de Transações: Lista paginada de todos os pagamentos processados pelo gateway.

Controle Manual: Aprovação ou rejeição de pagamentos com atualização imediata via API.

Filtros Inteligentes: Alternância rápida entre visualização total e apenas itens que aguardam ação.

Interface Responsiva: Construída com Tailwind CSS v4 para uma experiência fluida.

🛠 Tecnologias Utilizadas
React + Vite: Para uma interface rápida e reativa.

Tailwind CSS v4: Estilização moderna e otimizada.

Lucide React / Icons: Identificação visual de status e métodos.

Fetch API: Integração assíncrona com o backend em Go.

🏗 Arquitetura do Ecossistema
Este projeto faz parte de um sistema de microserviços:

Ecommerce API: Origem dos pedidos.

Gateway Payments: Engine de processamento (Go + RabbitMQ + MySQL).

Gateway Admin Panel (Este Repo): Interface de controle do Gateway.

🚦 Como Rodar o Projeto
Clone o repositório:

Bash
git clone https://github.com/SEU_USUARIO/gateway-admin-panel.git
Instale as dependências:

Bash
npm install
Configure a URL da API:
No arquivo App.jsx, certifique-se de que a variável GATEWAY_URL aponta para sua API Go (ex: http://localhost:8080).

Inicie o servidor de desenvolvimento:

Bash
npm run dev