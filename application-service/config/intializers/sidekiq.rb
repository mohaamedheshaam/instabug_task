require 'sidekiq'
require 'sidekiq/rabbitmq'

Sidekiq.configure_server do |config|
  config.redis = { url: 'redis://localhost:6379/0' }
end

Sidekiq.configure_client do |config|
  config.redis = { url: 'redis://localhost:6379/0' }
end

Sidekiq::Rabbitmq.configure(
  host: 'rabbitmq', 
  port: 5672, 
  vhost: '/',
  user: 'guest',
  password: 'guest'
)
