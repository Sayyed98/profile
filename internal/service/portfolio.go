package service

import (
	"context"
	"errors"
	"log"
	"net/mail"
	"strings"

	"github.com/mohdhujaifa/profile/internal/cache"
	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/mohdhujaifa/profile/internal/queue"
	"github.com/mohdhujaifa/profile/internal/worker"
)

var ErrInvalid = errors.New("invalid request")

type PortfolioStore interface {
	GetPortfolio(ctx context.Context) (model.Portfolio, error)
	InsertContact(ctx context.Context, m model.ContactMessage) (int64, error)
	UpdateContactStatus(ctx context.Context, id int64, status string) error
}

type PortfolioService struct {
	store     PortfolioStore
	cache     cache.PortfolioCache
	publisher queue.Publisher
	pool      *worker.Pool
}

func NewPortfolioService(store PortfolioStore, c cache.PortfolioCache, pub queue.Publisher, pool *worker.Pool) *PortfolioService {
	return &PortfolioService{store: store, cache: c, publisher: pub, pool: pool}
}

func (s *PortfolioService) GetPortfolio(ctx context.Context) (model.Portfolio, error) {
	if s.cache != nil {
		p, hit, err := s.cache.Get(ctx)
		if err != nil {
			log.Printf("cache get: %v", err)
		} else if hit {
			return p, nil
		}
	}
	p, err := s.store.GetPortfolio(ctx)
	if err != nil {
		return p, err
	}
	if s.cache != nil {
		if err := s.cache.Set(ctx, p); err != nil {
			log.Printf("cache set: %v", err)
		}
	}
	return p, nil
}

func (s *PortfolioService) SubmitContact(ctx context.Context, req model.ContactRequest) (int64, error) {
	req = sanitizeContact(req)
	if err := ValidateContact(req); err != nil {
		return 0, err
	}
	msg := model.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Body:    req.Message,
		Status:  "pending",
	}
	id, err := s.store.InsertContact(ctx, msg)
	if err != nil {
		return 0, err
	}
	msg.ID = id
	if s.pool != nil {
		ok := s.pool.Submit(func(jobCtx context.Context) {
			s.processContact(jobCtx, msg)
		})
		if !ok {
			s.processContact(ctx, msg)
		}
	} else {
		s.processContact(ctx, msg)
	}
	return id, nil
}

func (s *PortfolioService) processContact(ctx context.Context, msg model.ContactMessage) {
	status := "processed"
	if s.publisher != nil {
		if err := s.publisher.PublishContact(ctx, msg); err != nil {
			log.Printf("publish contact %d: %v", msg.ID, err)
			status = "failed"
		}
	}
	if err := s.store.UpdateContactStatus(ctx, msg.ID, status); err != nil {
		log.Printf("update contact status %d: %v", msg.ID, err)
	}
}

func sanitizeContact(req model.ContactRequest) model.ContactRequest {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Subject = strings.TrimSpace(req.Subject)
	req.Message = strings.TrimSpace(req.Message)
	return req
}

func ValidateContact(req model.ContactRequest) error {
	if req.Name == "" || len(req.Name) > 120 {
		return ErrInvalid
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return ErrInvalid
	}
	if req.Subject == "" || len(req.Subject) > 180 {
		return ErrInvalid
	}
	if req.Message == "" || len(req.Message) > 4000 {
		return ErrInvalid
	}
	return nil
}
