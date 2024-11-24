package handler

import (
    "net/http"
    "encoding/json"
    "go.uber.org/zap"
    "github.com/gorilla/mux"
    "chat-service/internal/service"
)

type ChatHandler struct {
    service *service.ChatService
    logger  *zap.Logger
}

func NewChatHandler(service *service.ChatService, logger *zap.Logger) *ChatHandler {
    return &ChatHandler{
        service: service,
        logger:  logger,
    }
}

// @Summary     Create a new chat
// @Description Creates a new chat for an application
// @Tags        chats
// @Accept      json
// @Produce     json
// @Param       token path string true "Application Token"
// @Success     201 {object} model.CreateChatResponse
// @Failure     500 {object} model.ErrorResponse
// @Router      /applications/{token}/chats [post]
func (h *ChatHandler) Create(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    applicationToken := vars["token"]

    h.logger.Info("creating new chat",
        zap.String("application_token", applicationToken))

    chat, err := h.service.CreateChat(r.Context(), applicationToken)
    if err != nil {
        h.logger.Error("failed to create chat",
            zap.Error(err),
            zap.String("application_token", applicationToken))
        respondWithError(w, http.StatusInternalServerError, "Failed to create chat")
        return
    }

    respondWithJSON(w, http.StatusCreated, map[string]interface{}{
        "chat_number": chat.Number,
    })
}

// @Summary     List all chats
// @Description Gets all chats for an application
// @Tags        chats
// @Accept      json
// @Produce     json
// @Param       token path string true "Application Token"
// @Success     200 {array} model.ChatResponse
// @Failure     400 {object} model.ErrorResponse
// @Failure     404 {object} model.ErrorResponse
// @Router      /applications/{token}/chats [get]
func (h *ChatHandler) ListChats(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    applicationToken := vars["token"]

    chats, err := h.service.ListChats(r.Context(), applicationToken)
    if err != nil {
        h.logger.Error("failed to list chats",
            zap.Error(err),
            zap.String("application_token", applicationToken))
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    response := make([]map[string]interface{}, len(chats))
    for i, chat := range chats {
        response[i] = map[string]interface{}{
            "number":         chat.Number,
            "messages_count": chat.MessagesCount,
            "created_at":     chat.CreatedAt,
        }
    }

    respondWithJSON(w, http.StatusOK, response)
}

// Helper functions for response handling
func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal Server Error"))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}