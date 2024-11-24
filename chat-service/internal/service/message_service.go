package service

import (
    "context"
    "fmt"
    "time"
    "strconv"
    
    "go.uber.org/zap"
    
    "chat-service/internal/model"
    "chat-service/internal/repository/mysql"
    "chat-service/internal/repository/redis"
    "chat-service/pkg/elasticsearch"
    "chat-service/pkg/rabbitmq"
)


type MessageService struct {
    messageRepo   *mysql.MessageRepository
    chatRepo     *mysql.ChatRepository
    sequenceRepo  *redis.SequenceRepository
    rabbitMQ      *rabbitmq.Client
    elasticSearch *elasticsearch.Client
    logger        *zap.Logger
}

func NewMessageService(
    messageRepo *mysql.MessageRepository,
    chatRepo *mysql.ChatRepository,
    sequenceRepo *redis.SequenceRepository,
    rabbitMQ *rabbitmq.Client,
    elasticSearch *elasticsearch.Client,
    logger *zap.Logger,
) *MessageService {
    return &MessageService{
        messageRepo:   messageRepo,
        chatRepo:     chatRepo,
        sequenceRepo:  sequenceRepo,
        rabbitMQ:      rabbitMQ,
        elasticSearch: elasticSearch,
        logger:        logger,
    }
}

func (s *MessageService) CreateMessage(ctx context.Context, applicationToken string, chatNumber string, body string) (*model.Message, error) {

    chatNum, err := strconv.Atoi(chatNumber)
    if err != nil {
        return nil, fmt.Errorf("invalid chat number: %w", err)
    }

    chat, err := s.chatRepo.GetByNumber(ctx, applicationToken, chatNum)
    if err != nil {
        return nil, fmt.Errorf("failed to get chat: %w", err)
    }
    if chat == nil {
        return nil, fmt.Errorf("chat not found")
    }

    number, err := s.sequenceRepo.NextMessageNumber(ctx, chat.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get next message number: %w", err)
    }

    message := &model.Message{
        ChatID:    chat.ID,
        Number:    number,
        Body:      body,
        CreatedAt: time.Now().UTC(),
    }

    if err := s.messageRepo.Create(ctx, message); err != nil {
        return nil, fmt.Errorf("failed to create message: %w", err)
    }

    go func() {
        messageJSON, err := message.ToJSON()
        if err != nil {
            s.logger.Error("failed to marshal message",
                zap.Error(err),
                zap.Uint64("message_id", message.ID))
            return
        }

        if err := s.elasticSearch.Index("messages", fmt.Sprintf("%d", message.ID), messageJSON); err != nil {
            s.logger.Error("failed to index message",
                zap.Error(err),
                zap.Uint64("message_id", message.ID))
        }
    }()

    go func() {
        if err := s.rabbitMQ.PublishMessageCreated(context.Background(), message); err != nil {
            s.logger.Error("failed to publish message created event",
                zap.Error(err),
                zap.Uint64("message_id", message.ID))
        }
    }()

    return message, nil
}

func (s *MessageService) ListMessages(ctx context.Context, applicationToken string, chatNumber string) ([]*model.Message, error) {
    chatNum, err := strconv.Atoi(chatNumber)
    if err != nil {
        return nil, fmt.Errorf("invalid chat number: %w", err)
    }

    chat, err := s.chatRepo.GetByNumber(ctx, applicationToken, chatNum)
    if err != nil {
        return nil, fmt.Errorf("failed to get chat: %w", err)
    }
    if chat == nil {
        return nil, fmt.Errorf("chat not found")
    }

    messages, err := s.messageRepo.ListByChat(ctx, chat.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to list messages: %w", err)
    }

    if messages == nil {
        return []*model.Message{}, nil
    }

    return messages, nil
}

func (s *MessageService) SearchMessages(ctx context.Context, applicationToken string, chatNumber string, query string) ([]*model.Message, error) {
    chatNum, err := strconv.Atoi(chatNumber)
    if err != nil {
        return nil, fmt.Errorf("invalid chat number: %w", err)
    }

    chat, err := s.chatRepo.GetByNumber(ctx, applicationToken, chatNum)
    if err != nil {
        return nil, fmt.Errorf("failed to get chat: %w", err)
    }
    if chat == nil {
        return nil, fmt.Errorf("chat not found")
    }

    searchQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {
                        "match": map[string]interface{}{
                            "body": query,
                        },
                    },
                    {
                        "term": map[string]interface{}{
                            "chat_id": chat.ID,
                        },
                    },
                },
            },
        },
        "sort": []map[string]interface{}{
            {
                "created_at": map[string]interface{}{
                    "order": "desc",
                },
            },
        },
    }

    messages, err := s.messageRepo.Search(ctx, searchQuery)
    if err != nil {
        s.logger.Error("failed to search messages",
            zap.Error(err),
            zap.String("application_token", applicationToken),
            zap.String("chat_number", chatNumber),
            zap.String("query", query))
        return nil, fmt.Errorf("failed to search messages: %w", err)
    }

    return messages, nil
}
