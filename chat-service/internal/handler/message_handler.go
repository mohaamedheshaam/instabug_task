package handler

import (
    "encoding/json"
    "net/http"
    "go.uber.org/zap"
    "github.com/gorilla/mux"
    
    // "chat-service/internal/model"
    "chat-service/internal/service"
    "chat-service/internal/util"
)

type MessageHandler struct {
    service *service.MessageService
    logger  *zap.Logger
}

func NewMessageHandler(service *service.MessageService, logger *zap.Logger) *MessageHandler {
    return &MessageHandler{
        service: service,
        logger:  logger,
    }
}

type CreateMessageRequest struct {
    Body string `json:"body"`
}


// @Summary     Create a message
// @Description Creates a new message in a specific chat
// @Tags        messages
// @Accept      json
// @Produce     json
// @Param       token  path string true "Application Token"
// @Param       number path int    true "Chat Number"
// @Param       body   body CreateMessageRequest true "Message Content"
// @Success     201 {object} model.CreateMessageResponse
// @Failure     400 {object} model.ErrorResponse
// @Failure     404 {object} model.ErrorResponse
// @Failure     500 {object} model.ErrorResponse
// @Router      /applications/{token}/chats/{number}/messages [post]
func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    applicationToken := vars["token"]
    chatNumber := vars["number"]

    var req CreateMessageRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("failed to decode request body",
            zap.Error(err))
        util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    message, err := h.service.CreateMessage(r.Context(), applicationToken, chatNumber, req.Body)
    if err != nil {
        h.logger.Error("failed to create message",
            zap.Error(err),
            zap.String("application_token", applicationToken),
            zap.String("chat_number", chatNumber))
        util.RespondWithError(w, http.StatusInternalServerError, "Failed to create message")
        return
    }

    util.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
        "message_number": message.Number,
    })
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
func (h *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    applicationToken := vars["token"]
    chatNumber := vars["number"]

    h.logger.Info("listing messages",
        zap.String("application_token", applicationToken),
        zap.String("chat_number", chatNumber))

    messages, err := h.service.ListMessages(r.Context(), applicationToken, chatNumber)
    if err != nil {
        h.logger.Error("failed to list messages",
            zap.Error(err),
            zap.String("application_token", applicationToken),
            zap.String("chat_number", chatNumber))
        util.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch messages")
        return
    }

    util.RespondWithJSON(w, http.StatusOK, messages)
}

// @Summary Search messages
// @Description Search for messages within a chat based on query text
// @Tags messages
// @Accept json
// @Produce json
// @Param token path string true "Application Token"
// @Param number path int true "Chat Number"
// @Param q query string true "Search Query"
// @Success 200 {array} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse "Search query is required"  
// @Failure 404 {object} model.ErrorResponse "Chat not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /applications/{token}/chats/{number}/messages/search [get]
func (h *MessageHandler) Search(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    applicationToken := vars["token"]
    chatNumber := vars["number"]
    query := r.URL.Query().Get("q")

    if query == "" {
        util.RespondWithError(w, http.StatusBadRequest, "Search query is required")
        return
    }

    h.logger.Info("searching messages",
        zap.String("application_token", applicationToken),
        zap.String("chat_number", chatNumber),
        zap.String("query", query))

    messages, err := h.service.SearchMessages(r.Context(), applicationToken, chatNumber, query)
    if err != nil {
        h.logger.Error("failed to search messages",
            zap.Error(err),
            zap.String("application_token", applicationToken),
            zap.String("chat_number", chatNumber),
            zap.String("query", query))
        util.RespondWithError(w, http.StatusInternalServerError, "Failed to search messages")
        return
    }

    util.RespondWithJSON(w, http.StatusOK, messages)
}