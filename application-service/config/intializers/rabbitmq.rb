Sneakers.configure(
  amqp: ENV['RABBITMQ_URL'],
  exchange: 'chat_events',
  exchange_type: :topic,
  durable: true
)