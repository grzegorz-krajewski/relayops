package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"relayops/apps/gateway-go/internal/metrics"
	"relayops/apps/gateway-go/internal/redisstream"
	"relayops/apps/gateway-go/internal/store"
)

type Handler struct {
	publisher *redisstream.Publisher
	store     *store.TaskStore
	stream    string
}

func NewHandler(publisher *redisstream.Publisher, store *store.TaskStore, stream string) *Handler {
	return &Handler{
		publisher: publisher,
		store:     store,
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

type TaskResponse struct {
	ID              string         `json:"id"`
	Type            string         `json:"type"`
	Status          string         `json:"status"`
	RawPayload      map[string]any `json:"raw_payload"`
	ResultPayload   map[string]any `json:"result_payload,omitempty"`
	TraceID         string         `json:"trace_id"`
	StreamMessageID string         `json:"stream_message_id,omitempty"`
	ErrorMessage    string         `json:"error_message,omitempty"`
	CreatedAt       string         `json:"created_at"`
	ProcessedAt     string         `json:"processed_at,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/tasks", h.handleTasks)
	mux.HandleFunc("/api/v1/tasks/", h.handleGetTask)
}

func (h *Handler) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.handleCreateTask(w, r)
		return
	}
	writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
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

	taskID := uuid.NewString()
	traceID := uuid.NewString()

	rawPayloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "payload could not be serialized"})
		return
	}

	messageID, err := h.publisher.PublishTask(
		r.Context(),
		taskID,
		req.Type,
		string(rawPayloadBytes),
		traceID,
	)
	if err != nil {
		metrics.TaskPublishErrorsTotal.Inc()
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to publish task"})
		return
	}

	if err := h.store.InsertTask(
		r.Context(),
		taskID,
		req.Type,
		"accepted",
		req.Payload,
		traceID,
		messageID,
	); err != nil {
		metrics.TaskPersistErrorsTotal.Inc()
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to persist task"})
		return
	}

	writeJSON(w, http.StatusAccepted, CreateTaskResponse{
		TaskID:    taskID,
		Status:    "accepted",
		Stream:    h.stream,
		TraceID:   traceID,
		MessageID: messageID,
	})

	metrics.TasksCreatedTotal.Inc()
}

func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	taskID := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "task id is required"})
		return
	}

	rec, err := h.store.GetTaskByID(r.Context(), taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch task"})
		return
	}

	resp := TaskResponse{
		ID:        rec.ID,
		Type:      rec.Type,
		Status:    rec.Status,
		TraceID:   rec.TraceID,
		CreatedAt: rec.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}

	if rec.StreamMessageID.Valid {
		resp.StreamMessageID = rec.StreamMessageID.String
	}
	if rec.ErrorMessage.Valid {
		resp.ErrorMessage = rec.ErrorMessage.String
	}
	if rec.ProcessedAt.Valid {
		resp.ProcessedAt = rec.ProcessedAt.Time.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}

	_ = json.Unmarshal(rec.RawPayload, &resp.RawPayload)
	if len(rec.ResultPayload) > 0 {
		_ = json.Unmarshal(rec.ResultPayload, &resp.ResultPayload)
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
