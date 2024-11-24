package model

import (
    "time"
)

type Chat struct {
    ID            uint64    `json:"id"`
    ApplicationID string    `json:"application_id"`
    Number        int       `json:"number"`
    MessagesCount int       `json:"messages_count"`
    CreatedAt     time.Time `json:"created_at"`
}
