class ChatConsumer
    include Sneakers::Worker
  
    from_queue 'chat_created'
  
    def work(message)
      payload = JSON.parse(message)
      ChatCounterWorker.perform_async(payload['application_token'])
      ack!
    end
  end
  