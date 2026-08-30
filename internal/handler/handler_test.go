package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohdhujaifa/profile/internal/cache"
	"github.com/mohdhujaifa/profile/internal/handler"
	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/mohdhujaifa/profile/internal/queue"
	"github.com/mohdhujaifa/profile/internal/service"
	"github.com/stretchr/testify/require"
)

type stubStore struct{}

func (stubStore) GetPortfolio(context.Context) (model.Portfolio, error) {
	return model.Portfolio{Profile: model.Profile{FullName: "Mohd Hujaifa", Title: "Software Engineer"}}, nil
}
func (stubStore) InsertContact(context.Context, model.ContactMessage) (int64, error) { return 9, nil }
func (stubStore) UpdateContactStatus(context.Context, int64, string) error           { return nil }

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := service.NewPortfolioService(stubStore{}, cache.NewMemory(), queue.NoopPublisher{}, nil)
	s := handler.New(svc)
	r := gin.New()
	s.Register(r, []string{"http://localhost:5173"})
	return r
}

func TestHealth(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetPortfolio(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/portfolio", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var p model.Portfolio
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &p))
	require.Equal(t, "Mohd Hujaifa", p.Profile.FullName)
}

func TestPostContactRejectsInvalid(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/contact", bytes.NewBufferString(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostContactAccepted(t *testing.T) {
	r := setupRouter()
	body := `{"name":"Alex","email":"alex@example.com","subject":"Hello","message":"Interested in working together"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/contact", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusAccepted, w.Code)
}

func TestCORSPreflight(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/portfolio", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestAllowedOrigin(t *testing.T) {
	require.True(t, handler.AllowedOrigin([]string{"http://localhost:5173"}, "http://localhost:5173"))
	require.False(t, handler.AllowedOrigin([]string{"http://localhost:5173"}, "http://evil.example"))
}
