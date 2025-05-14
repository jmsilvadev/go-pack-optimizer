package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/optimizer"
)

// HandlerInterface defines the HTTP handler contract for pack optimizer endpoints.
type HandlerInterface interface {
	HealthHandler(w http.ResponseWriter, r *http.Request)
	CalculateOrder(w http.ResponseWriter, r *http.Request)
	GetPacks(w http.ResponseWriter, r *http.Request)
	PostPacks(w http.ResponseWriter, r *http.Request)
	DeletePacks(w http.ResponseWriter, r *http.Request)
	NotFoundHandler(w http.ResponseWriter, r *http.Request)
}

// Handler implements HTTP endpoints for managing and calculating packaging sizes.
type Handler struct {
	optimizer optimizer.OptimizerInterface
}

// Response defines a generic response structure for all endpoints.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// New creates a new Handler instance with the given optimizer implementation.
func New(opt optimizer.OptimizerInterface) *Handler {
	return &Handler{
		optimizer: opt,
	}
}

// GetPacks handles GET /v1/packs
// Returns the list of available pack sizes from the optimizer.
func (h *Handler) GetPacks(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	sizes, err := h.optimizer.GetAllSizes()
	if err != nil {
		response.Message = "no sizes found"
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	response.Data = sizes
	writeJSONResponse(w, http.StatusOK, response)
}

// PostPacks handles POST /v1/packs
// Adds a new pack size to the system via the optimizer.
func (h *Handler) PostPacks(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Size int `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Size <= 0 {
		http.Error(w, "Invalid or missing pack size", http.StatusBadRequest)
		return
	}

	if err := h.optimizer.AddSize(req.Size); err != nil {
		http.Error(w, "Failed to add pack size", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "size added succesfully",
	}
	writeJSONResponse(w, http.StatusCreated, response)
}

// DeletePacks handles DELETE /v1/packs/{size}
// Removes a specific pack size based on the size provided in the URL path.
func (h *Handler) DeletePacks(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid pack size", http.StatusBadRequest)
		return
	}

	if err := h.optimizer.RemoveSize(size); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := Response{}
	writeJSONResponse(w, http.StatusNoContent, response)
}

// CalculateOrder handles POST /v1/order
// Calculates the optimal set of packs to fulfill a given quantity.
func (h *Handler) CalculateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemsOrdered int `json:"items_ordered"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ItemsOrdered <= 0 {
		http.Error(w, "items_ordered must be greater than 0", http.StatusBadRequest)
		return
	}

	result := h.optimizer.Calculate(req.ItemsOrdered)

	resp := map[string]interface{}{
		"packs":       result.PacksUsed,
		"total_items": result.TotalItems,
		"total_packs": result.TotalPacks,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// NotFoundHandler handles requests to undefined routes.
func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "route not found",
	}
	writeJSONResponse(w, http.StatusNotFound, response)
}

// HealthHandler handles GET /health
// Returns a simple success response for health checks.
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	writeJSONResponse(w, http.StatusOK, response)
}

// writeJSONResponse encodes and writes a JSON response with the given status code.
func writeJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	response.Status = "error"
	if statusCode < 400 {
		response.Status = "success"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
