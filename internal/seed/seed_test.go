package seed_test

import (
	"context"
	"testing"

	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/mohdhujaifa/profile/internal/seed"
	"github.com/stretchr/testify/require"
)

type memStore struct {
	has      bool
	profile  model.Profile
	skills   int
	exps     int
	projects int
	edu      int
}

func (m *memStore) HasProfile(context.Context) (bool, error) { return m.has, nil }
func (m *memStore) InsertProfile(_ context.Context, p model.Profile) error {
	m.profile = p
	m.has = true
	return nil
}
func (m *memStore) InsertSkill(context.Context, string, string, int) error {
	m.skills++
	return nil
}
func (m *memStore) InsertExperience(context.Context, model.Experience) (int64, error) {
	m.exps++
	return int64(m.exps), nil
}
func (m *memStore) InsertProject(context.Context, model.Project) error {
	m.projects++
	return nil
}
func (m *memStore) InsertEducation(context.Context, model.Education) error {
	m.edu++
	return nil
}

func TestResumeIfEmpty(t *testing.T) {
	s := &memStore{}
	require.NoError(t, seed.ResumeIfEmpty(context.Background(), s))
	require.Equal(t, "Mohd Hujaifa", s.profile.FullName)
	require.Greater(t, s.skills, 10)
	require.Equal(t, 3, s.exps)
	require.Equal(t, 3, s.projects)
	require.Equal(t, 2, s.edu)

	skillsBefore := s.skills
	require.NoError(t, seed.ResumeIfEmpty(context.Background(), s))
	require.Equal(t, skillsBefore, s.skills)
}
