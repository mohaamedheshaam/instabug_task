package redis

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
)

type SequenceRepository struct {
    client *redis.Client
}

func NewSequenceRepository(client *redis.Client) *SequenceRepository {
    return &SequenceRepository{client: client}
}

func (r *SequenceRepository) NextChatNumber(ctx context.Context, applicationID string) (int, error) {
    key := fmt.Sprintf("app:%s:chat_seq", applicationID)
    fmt.Printf("Requesting next chat number for application: %s, key: %s\n", applicationID, key)
    return r.getNextSequence(ctx, key)
}

func (r *SequenceRepository) NextMessageNumber(ctx context.Context, chatID uint64) (int, error) {
    key := fmt.Sprintf("chat:%d:msg_seq", chatID)
    fmt.Printf("Requesting next message number for chatID: %d, key: %s\n", chatID, key)
    return r.getNextSequence(ctx, key)
}

func (r *SequenceRepository) getNextSequence(ctx context.Context, key string) (int, error) {
    fmt.Printf("Incrementing Redis key: %s\n", key)
    val, err := r.client.Incr(ctx, key).Result()
    if err != nil {
        return 0, fmt.Errorf("failed to increment sequence: %w", err)
    }
    fmt.Printf("Successfully incremented key: %s, new value: %d\n", key, val)
    return int(val), nil
}
