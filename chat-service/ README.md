# Chat Service

Fast, scalable chat management service built with Go. Handles chat creation, message management, and real-time search functionality.

## 🚀 Features

* Real-time message search with Elasticsearch
* Sequential chat and message numbering
* Event-driven architecture using RabbitMQ
* Atomic counters with Redis
* RESTful API with Swagger documentation

## 🔧 Configuration

Environment variables:
```bash
# MySQL
MYSQL_HOST=db
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DATABASE=chat

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# Elasticsearch
ELASTICSEARCH_URL=http://elasticsearch:9200
```

## 🛣️ API Routes

### Chats
- `POST /api/applications/{token}/chats` - Create chat

### Messages
- `POST /api/applications/{token}/chats/{number}/messages` - Create message
- `GET /api/applications/{token}/chats/{number}/messages` - List messages
- `GET /api/applications/{token}/chats/{number}/messages/search` - Search messages

## 📚 Database Schema

```sql

CREATE TABLE IF NOT EXISTS chats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id VARCHAR(255) NOT NULL,
    number INT NOT NULL,
    messages_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    UNIQUE KEY unique_app_number (application_id, number),
    FOREIGN KEY (application_id) REFERENCES applications(token)
);

CREATE TABLE IF NOT EXISTS messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    chat_id BIGINT UNSIGNED NOT NULL,
    number INT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    UNIQUE KEY unique_chat_number (chat_id, number)
);
```

## 🏗️ Architecture

- **Redis**: Atomic sequence generation
- **RabbitMQ**: Event publishing
- **Elasticsearch**: Message searching
- **MySQL**: Data persistence

## 📖 API Documentation
Swagger UI available at: `http://localhost:8080/swagger/index.html`

## 🐳 Docker

Run with Docker Compose:
```bash
docker-compose up go-service
```