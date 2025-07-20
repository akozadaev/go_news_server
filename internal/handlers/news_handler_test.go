package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_news_server/internal/models"
	"go_news_server/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNewsService - мок для NewsService
type MockNewsService struct {
	mock.Mock
}

func (m *MockNewsService) UpdateNews(ctx context.Context, news *models.News, categories []int64) error {
	args := m.Called(ctx, news, categories)
	return args.Error(0)
}

func (m *MockNewsService) GetNewsList(ctx context.Context, limit, offset int) ([]models.News, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.News), args.Error(1)
}

// Убедимся, что MockNewsService реализует интерфейс
var _ services.NewsServiceInterface = (*MockNewsService)(nil)

func TestEditNewsHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockNewsService)
	handler := &NewsHandlers{Service: mockService}

	app.Post("/edit/:Id", handler.EditNewsHandler)

	tests := []struct {
		name           string
		id             string
		requestBody    map[string]interface{}
		expectedStatus int
		setupMock      func()
	}{
		{
			name: "Valid update request",
			id:   "1",
			requestBody: map[string]interface{}{
				"Id":         1,
				"Title":      "Updated Title",
				"Content":    "Updated Content",
				"Categories": []int64{1, 2, 3},
			},
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mockService.On("UpdateNews", mock.Anything, mock.Anything, []int64{1, 2, 3}).Return(nil)
			},
		},
		{
			name:           "Invalid ID",
			id:             "invalid",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/edit/"+tt.id, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)
			defer func() {
				if err := resp.Body.Close(); err != nil {
				}
			}()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetNewsList(t *testing.T) {
	app := fiber.New()
	mockService := new(MockNewsService)
	handler := &NewsHandlers{Service: mockService}

	app.Get("/list", handler.GetNewsList)

	expectedNews := []models.News{
		{
			Id:         1,
			Title:      "Test News",
			Content:    "Test Content",
			Categories: []int64{1, 2},
		},
	}

	mockService.On("GetNewsList", mock.Anything, 10, 0).Return(expectedNews, nil)

	req := httptest.NewRequest("GET", "/list", nil)
	resp, _ := app.Test(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
		}
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.True(t, result["Success"].(bool))
	assert.NotNil(t, result["News"])

	mockService.AssertExpectations(t)
}
