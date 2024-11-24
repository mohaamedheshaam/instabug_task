package mysql

import (
    "context"
    "database/sql" 
    "encoding/json"
    "fmt"
    
    "chat-service/internal/model"
    "chat-service/pkg/elasticsearch"
)

type MessageRepository struct {
    db *sql.DB
    es *elasticsearch.Client
}

func NewMessageRepository(db *sql.DB, es *elasticsearch.Client) *MessageRepository {
    return &MessageRepository{
        db: db,
        es: es,
    }
}


// @Summary     Create a message
// @Description Creates a new message in a specific chat
// @Tags        messages
// @Accept      json
// @Produce     json
// @Param       token  path string true "Application Token"
// @Param       number path int    true "Chat Number"
// @Param       body   body model.CreateMessageRequest true "Message Content"
// @Success     201 {object} model.CreateMessageResponse
// @Failure     400 {object} model.ErrorResponse
// @Failure     404 {object} model.ErrorResponse
// @Failure     500 {object} model.ErrorResponse
// @Router      /applications/{token}/chats/{number}/messages [post]
func (r *MessageRepository) Create(ctx context.Context, message *model.Message) error {
    query := `
        INSERT INTO messages (chat_id, number, body, created_at)
        VALUES (?, ?, ?, ?)
    `
    
    result, err := r.db.ExecContext(ctx, query,
        message.ChatID,
        message.Number,
        message.Body,
        message.CreatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to insert message: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert id: %w", err)
    }

    message.ID = uint64(id)
    
    if err := r.es.Index("messages", fmt.Sprintf("%d", message.ID), message); err != nil {
        return fmt.Errorf("failed to index message: %w", err)
    }
    
    return nil
}


// @Summary     List messages
// @Description Retrieves all messages from a specific chat
// @Tags        messages
// @Produce     json
// @Param       token  path string true "Application Token"
// @Param       number path int    true "Chat Number"
// @Success     200 {array} model.MessageResponse
// @Failure     404 {object} model.ErrorResponse
// @Failure     500 {object} model.ErrorResponse
// @Router      /applications/{token}/chats/{number}/messages [get]
func (r *MessageRepository) ListByChat(ctx context.Context, chatID uint64) ([]*model.Message, error) {
    query := `
        SELECT id, chat_id, number, body, created_at
        FROM messages
        WHERE chat_id = ?
        ORDER BY number ASC
    `
    
    rows, err := r.db.QueryContext(ctx, query, chatID)
    if err != nil {
        return nil, fmt.Errorf("failed to query messages: %w", err)
    }
    defer rows.Close()

    var messages []*model.Message
    for rows.Next() {
        msg := &model.Message{}
        err := rows.Scan(
            &msg.ID,
            &msg.ChatID,
            &msg.Number,
            &msg.Body,
            &msg.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan message: %w", err)
        }
        messages = append(messages, msg)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating messages: %w", err)
    }
    
    return messages, nil
}


// @Summary     Search messages
// @Description Searches for messages within a specific chat
// @Tags        messages
// @Produce     json
// @Param       token  path  string true "Application Token"
// @Param       number path  int    true "Chat Number"
// @Param       instabug      query string true "Search Query"
// @Success     200 {array} model.MessageResponse
// @Failure     400 {object} model.ErrorResponse
// @Failure     404 {object} model.ErrorResponse
// @Failure     500 {object} model.ErrorResponse
// @Router      /applications/{token}/chats/{number}/messages/search [get]

func (r *MessageRepository) Search(ctx context.Context, query map[string]interface{}) ([]*model.Message, error) {
    searchResults, err := r.es.Search("messages", query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute search: %w", err)
    }

    var searchResponse struct {
        Hits struct {
            Hits []struct {
                Source *model.Message `json:"_source"`
            } `json:"hits"`
        } `json:"hits"`
    }

    if err := json.Unmarshal(searchResults, &searchResponse); err != nil {
        return nil, fmt.Errorf("failed to parse search results: %w", err)
    }

    messages := make([]*model.Message, len(searchResponse.Hits.Hits))
    for i, hit := range searchResponse.Hits.Hits {
        messages[i] = hit.Source
    }

    return messages, nil
}