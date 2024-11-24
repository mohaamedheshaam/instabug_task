class MessageCounterWorker
    include Sidekiq::Worker
  
    def perform(chat_id)
      Chat.where(id: chat_id)
          .update_all('messages_count = messages_count + 1')
    end
  end
  