package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"relayops/apps/gateway-go/internal/redisstream"
)

type Handler struct {
	publisher *redisstream.Publisher
	stream    string
}

func NewHandler(publisher *redisstream.Publisher, stream string) *Handler {
	return &Handler{
		publisher: publisher,
		stream:    stream,
	}
}

type CreateTaskRequest struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

type CreateTaskResponse struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	Stream    string `json:"stream"`
	TraceID   string `json:"trace_id"`
	MessageID string `json:"message_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/tasks", h.handleCreateTask)
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return
	}

	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "type is required"})
		return
	}

	if req.Payload == nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "payload is required"})
		return
	}

	rawPayloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "payload could not be serialized"})
		return
	}

	taskID := uuid.NewString()
	traceID := uuid.NewString()

	messageID, err := h.publisher.PublishTask(
		r.Context(),
		taskID,
		req.Type,
		string(rawPayloadBytes),
		traceID,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to publish task"})
		return
	}

	writeJSON(w, http.StatusAccepted, CreateTaskResponse{
		TaskID:    taskID,
		Status:    "accepted",
		Stream:    h.stream,
		TraceID:   traceID,
		MessageID: messageID,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
