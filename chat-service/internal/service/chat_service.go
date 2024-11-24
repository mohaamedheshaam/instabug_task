package service

import (
    "context"
    "fmt"
    "time"
    
    "go.uber.org/zap"
    
    "chat-service/internal/model"
    "chat-service/internal/repository/mysql"
    "chat-service/internal/repository/redis"
    "chat-service/pkg/rabbitmq"
)

type ChatService struct {
    chatRepo     *mysql.ChatRepository
    sequenceRepo *redis.SequenceRepository
    rabbitMQ     *rabbitmq.Client
    logger       *zap.Logger
}

func NewChatService(
    chatRepo *mysql.ChatRepository,
    sequenceRepo *redis.SequenceRepository,
    rabbitMQ *rabbitmq.Client,
    logger *zap.Logger,
) *ChatService {
    return &ChatService{
        chatRepo:     chatRepo,
        sequenceRepo: sequenceRepo,
        rabbitMQ:     rabbitMQ,
        logger:       logger,
    }
}

func (s *ChatService) CreateChat(ctx context.Context, applicationID string) (*model.Chat, error) {
    number, err := s.sequenceRepo.NextChatNumber(ctx, applicationID)
    if err != nil {
        return nil, fmt.Errorf("failed to get next chat number: %w", err)
    }

    chat := &model.Chat{
        ApplicationID: applicationID,
        Number:       number,
        CreatedAt:    time.Now().UTC(),
    }

    if err := s.chatRepo.Create(ctx, chat); err != nil {
        return nil, fmt.Errorf("failed to create chat: %w", err)
    }

    go func() {
        if err := s.rabbitMQ.PublishChatCreated(context.Background(), chat); err != nil {
            s.logger.Error("failed to publish chat created event",
                zap.Error(err),
                zap.String("application_id", applicationID),
                zap.Int("number", number))
        }
    }()

    return chat, nil
}

func (s *ChatService) ListChats(ctx context.Context, applicationToken string) ([]*model.Chat, error) {
    chats, err := s.chatRepo.ListByApplication(ctx, applicationToken)
    if err != nil {
        return nil, fmt.Errorf("failed to list chats: %w", err)
    }
    
    return chats, nil
}