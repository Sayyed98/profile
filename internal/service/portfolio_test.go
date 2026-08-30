package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mohdhujaifa/profile/internal/cache"
	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/mohdhujaifa/profile/internal/service"
	"github.com/mohdhujaifa/profile/internal/worker"
	"github.com/stretchr/testify/require"
)

type fakeStore struct {
	mu        sync.Mutex
	portfolio model.Portfolio
	contacts  []model.ContactMessage
	statuses  map[int64]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		portfolio: model.Portfolio{
			Profile: model.Profile{FullName: "Mohd Hujaifa", Title: "Software Engineer"},
		},
		statuses: map[int64]string{},
	}
}

func (f *fakeStore) GetPortfolio(context.Context) (model.Portfolio, error) {
	return f.portfolio, nil
}

func (f *fakeStore) InsertContact(_ context.Context, m model.ContactMessage) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	m.ID = int64(len(f.contacts) + 1)
	f.contacts = append(f.contacts, m)
	return m.ID, nil
}

func (f *fakeStore) UpdateContactStatus(_ context.Context, id int64, status string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.statuses[id] = status
	return nil
}

type recordingPublisher struct {
	mu   sync.Mutex
	msgs []model.ContactMessage
}

func (r *recordingPublisher) PublishContact(_ context.Context, msg model.ContactMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.msgs = append(r.msgs, msg)
	return nil
}

func (r *recordingPublisher) Close() error { return nil }

func TestGetPortfolioUsesCache(t *testing.T) {
	store := newFakeStore()
	mem := cache.NewMemory()
	svc := service.NewPortfolioService(store, mem, nil, nil)

	first, err := svc.GetPortfolio(context.Background())
	require.NoError(t, err)
	require.Equal(t, "Mohd Hujaifa", first.Profile.FullName)

	store.portfolio.Profile.FullName = "Changed"
	second, err := svc.GetPortfolio(context.Background())
	require.NoError(t, err)
	require.Equal(t, "Mohd Hujaifa", second.Profile.FullName)
}

func TestSubmitContactValidatesAndPublishes(t *testing.T) {
	store := newFakeStore()
	pub := &recordingPublisher{}
	pool := worker.NewPool(context.Background(), 2)
	defer pool.Stop()
	svc := service.NewPortfolioService(store, cache.NewMemory(), pub, pool)

	_, err := svc.SubmitContact(context.Background(), model.ContactRequest{
		Name: "A", Email: "bad", Subject: "Hi", Message: "Hello",
	})
	require.ErrorIs(t, err, service.ErrInvalid)

	id, err := svc.SubmitContact(context.Background(), model.ContactRequest{
		Name: "Recruiter", Email: "recruiter@example.com", Subject: "Role", Message: "Let's talk",
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), id)

	require.Eventually(t, func() bool {
		pub.mu.Lock()
		defer pub.mu.Unlock()
		return len(pub.msgs) == 1
	}, time.Second, 10*time.Millisecond)
}

func TestValidateContact(t *testing.T) {
	require.Error(t, service.ValidateContact(model.ContactRequest{}))
	require.NoError(t, service.ValidateContact(model.ContactRequest{
		Name: "Hujaifa", Email: "huzaifaanis40@gmail.com", Subject: "Hello", Message: "Nice portfolio",
	}))
}
