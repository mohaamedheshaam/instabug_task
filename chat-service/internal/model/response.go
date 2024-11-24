package model

import "time"

type CreateChatResponse struct {
    ChatNumber int `json:"chat_number" example:"1"`
}

type CreateMessageRequest struct {
    Body string `json:"body" example:"Welcome to instabug!!" binding:"required"`
}

type CreateMessageResponse struct {
    MessageNumber int `json:"message_number" example:"1"`
}

type MessageResponse struct {
    ID        uint64    `json:"id" example:"1"`
    ChatID    uint64    `json:"chat_id" example:"1"`
    Number    int       `json:"number" example:"1"`
    Body      string    `json:"body" example:"Welcome to instabug!!"`
    CreatedAt time.Time `json:"created_at" example:"2024-11-19T20:00:00Z"`
}

type ErrorResponse struct {
    Error string `json:"error" example:"Error message"`
}
