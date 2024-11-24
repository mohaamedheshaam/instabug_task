class MessageConsumer
    include Sneakers::Worker
  
    from_queue 'message_created'
  
    def work(message)
      payload = JSON.parse(message)
      MessageCounterWorker.perform_async(payload['chat_id'])
      ack!
    end
  end
  