# Application Management Service

Application management and event processing service built with Ruby on Rails. Handles application creation and counter management.

## 🚀 Features

* Application CRUD operations
* Asynchronous counter updates
* Event processing with RabbitMQ
* Background job processing with Sidekiq

## 🔧 Configuration

Environment variables:
```bash
# Database
DATABASE_URL=mysql2://root:password@db:3306/chat

# Redis
REDIS_URL=redis://redis:6379/0

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672

# Rails
RAILS_ENV=development
```

## 🛣️ API Routes

### Applications
- `POST /api/v1/applications` - Create application
- `GET /api/v1/applications` - List applications
- `GET /api/v1/applications/:token` - Get application

## 📚 Database Schema

```sql
CREATE TABLE applications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    chats_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE KEY unique_token (token)
);
```

## 👷 Workers

- **ChatCounterWorker**: Updates chat counts
- **MessageCounterWorker**: Updates message counts
## 🏗️ Architecture

- **Redis**: Sidekiq backend
- **RabbitMQ**: Event consumption
- **MySQL**: Data storage

## 📊 Monitoring

- RabbitMQ management: `http://localhost:15672`
  - Username: `guest`
  - Password: `guest`

## 🐳 Docker

Run with Docker Compose:
```bash
docker-compose up rails-service
```