Chat System
A distributed chat system built with Go and Ruby on Rails, featuring real-time search and asynchronous event processing.
🏗️ Architecture

Go Service: Chat and message management
Rails Service: Application management and counter updates
MySQL: Primary data store
Redis: Sequence generation and caching
RabbitMQ: Event messaging
Elasticsearch: Message searching

🚀 System Overview
Show Image

Go Service: localhost:8080
Rails Service: localhost:3000
RabbitMQ Dashboard: localhost:15672
Elasticsearch: localhost:9200


Clone the repository 

Start all services:

docker-compose up

Wait for all services to be ready:


MySQL will initialize with required tables
Elasticsearch will start and be ready for indexing
RabbitMQ management interface will be accessible
Both Go and Rails services will start
