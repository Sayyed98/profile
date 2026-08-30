package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/mohdhujaifa/profile/internal/service"
)

type Server struct {
	svc *service.PortfolioService
}

func New(svc *service.PortfolioService) *Server {
	return &Server{svc: svc}
}

func (s *Server) Register(r *gin.Engine, corsOrigins []string) {
	r.Use(corsMiddleware(corsOrigins))
	r.GET("/healthz", s.health)
	api := r.Group("/api/v1")
	{
		api.GET("/portfolio", s.getPortfolio)
		api.POST("/contact", s.postContact)
	}
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) getPortfolio(c *gin.Context) {
	p, err := s.svc.GetPortfolio(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load portfolio"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (s *Server) postContact(c *gin.Context) {
	var req model.ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	id, err := s.svc.SubmitContact(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contact payload"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit message"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": id, "status": "accepted"})
}

func corsMiddleware(origins []string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, o := range origins {
		allowed[o] = struct{}{}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok || len(allowed) == 0 {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Headers", "Content-Type")
				c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			}
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func AllowedOrigin(origins []string, origin string) bool {
	origin = strings.TrimSpace(origin)
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}
