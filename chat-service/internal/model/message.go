
package model

import (
    "encoding/json"
    "time"
)

type Message struct {
    ID        uint64    `json:"id"`
    ChatID    uint64    `json:"chat_id"`
    Number    int       `json:"number"`
    Body      string    `json:"body"`
    CreatedAt time.Time `json:"created_at"`
}

func (m *Message) ToJSON() ([]byte, error) {
    return json.Marshal(m)
}
