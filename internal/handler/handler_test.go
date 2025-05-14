package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmsilvadev/go-pack-optimizer/internal/handler/mocks"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOptimizer := mocks.NewMockOptimizerInterface(ctrl)

	handler := New(mockOptimizer)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.HealthHandler(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCalculateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOptimizer := mocks.NewMockOptimizerInterface(ctrl)

	handler := New(mockOptimizer)

	mockOptimizer.EXPECT().Calculate(501).Return(&optimizer.OptimizationResult{
		PacksUsed:  []int{500, 250},
		TotalItems: 750,
		TotalPacks: 2,
	}).Times(1)

	body := map[string]int{"items_ordered": 501}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/order", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	handler.CalculateOrder(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result optimizer.OptimizationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, 750, result.TotalItems)
	assert.Equal(t, 2, result.TotalPacks)

	req = httptest.NewRequest("POST", "/v1/order", nil)
	rr = httptest.NewRecorder()

	handler.CalculateOrder(rr, req)

	resp = rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetPacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOptimizer := mocks.NewMockOptimizerInterface(ctrl)

	handler := New(mockOptimizer)

	mockOptimizer.EXPECT().GetAllSizes().Return([]int{250, 500, 1000}, nil).Times(1)

	req := httptest.NewRequest("GET", "/v1/packs", nil)
	rr := httptest.NewRecorder()

	handler.GetPacks(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var payload Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, float64(250), payload.Data.([]interface{})[0])
	assert.Equal(t, float64(500), payload.Data.([]interface{})[1])
	assert.Equal(t, float64(1000), payload.Data.([]interface{})[2])
}

func TestPostPacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOptimizer := mocks.NewMockOptimizerInterface(ctrl)

	handler := New(mockOptimizer)

	mockOptimizer.EXPECT().AddSize(1500).Return(nil).Times(1)

	body := map[string]int{"size": 1500}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/v1/packs", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	handler.PostPacks(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	req = httptest.NewRequest("POST", "/v1/packs", nil)
	rr = httptest.NewRecorder()

	handler.PostPacks(rr, req)

	resp = rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeletePacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOptimizer := mocks.NewMockOptimizerInterface(ctrl)

	handler := New(mockOptimizer)

	mockOptimizer.EXPECT().RemoveSize(500).Return(nil).Times(1)

	req := httptest.NewRequest("DELETE", "/v1/packs/500", nil)
	rr := httptest.NewRecorder()

	handler.DeletePacks(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	req = httptest.NewRequest("DELETE", "/v1/packs/", nil)
	rr = httptest.NewRecorder()

	handler.DeletePacks(rr, req)

	resp = rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	req = httptest.NewRequest("DELETE", "/v1/packs/a", nil)
	rr = httptest.NewRecorder()

	handler.DeletePacks(rr, req)

	resp = rr.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestNotFoundHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOptimizer := mocks.NewMockOptimizerInterface(ctrl)

	handler := New(mockOptimizer)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.NotFoundHandler(rr, req)

	resp := rr.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestNewRouter(t *testing.T) {
	_, err := NewRouter(nil)
	assert.Error(t, err)

	_, err = NewRouter(&Handler{})
	assert.NoError(t, err)
}
