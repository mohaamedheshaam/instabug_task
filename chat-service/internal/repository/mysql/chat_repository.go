package mysql

import (
    "context"
    "database/sql"
    "fmt"
    "chat-service/internal/model"
)

type ChatRepository struct {
    db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
    return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *model.Chat) error {
    query := `
        INSERT INTO chats (application_id, number, messages_count, created_at)
        VALUES (?, ?, ?, ?)
    `
    fmt.Printf("Creating chat with application_id: %s, number: %d, messages_count: %d, created_at: %s\n", 
    chat.ApplicationID, chat.Number, chat.MessagesCount, chat.CreatedAt) 
    result, err := r.db.ExecContext(ctx, query,
        chat.ApplicationID,
        chat.Number,
        chat.MessagesCount,
        chat.CreatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to insert chat: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert id: %w", err)
    }

    chat.ID = uint64(id)
    return nil
}

func (r *ChatRepository) GetByNumber(ctx context.Context, applicationID string, number int) (*model.Chat, error) {
    query := `
        SELECT id, application_id, number, messages_count, created_at
        FROM chats
        WHERE application_id = ? AND number = ?
    `
    
    chat := &model.Chat{}
    err := r.db.QueryRowContext(ctx, query, applicationID, number).Scan(
        &chat.ID,
        &chat.ApplicationID,
        &chat.Number,
        &chat.MessagesCount,
        &chat.CreatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to query chat: %w", err)
    }
    
    return chat, nil
}


func (r *ChatRepository) ListByApplication(ctx context.Context, applicationToken string) ([]*model.Chat, error) {
    query := `
        SELECT c.id, c.application_id, c.number, c.messages_count, c.created_at
        FROM chats c
        JOIN applications a ON c.application_id = a.id
        WHERE a.token = ?
        ORDER BY c.number ASC
    `
    
    rows, err := r.db.QueryContext(ctx, query, applicationToken)
    if err != nil {
        return nil, fmt.Errorf("failed to query chats: %w", err)
    }
    defer rows.Close()

    var chats []*model.Chat
    for rows.Next() {
        chat := &model.Chat{}
        err := rows.Scan(
            &chat.ID,
            &chat.ApplicationID,
            &chat.Number,
            &chat.MessagesCount,
            &chat.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan chat: %w", err)
        }
        chats = append(chats, chat)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating chats: %w", err)
    }
    
    return chats, nil
}
