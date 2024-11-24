class ChatCounterWorker
    include Sidekiq::Worker
  
    def perform(application_token)
      Application.where(token: application_token)
                 .update_all('chats_count = chats_count + 1')
    end
  end
  