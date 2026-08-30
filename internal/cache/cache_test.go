package cache_test

import (
	"context"
	"testing"

	"github.com/mohdhujaifa/profile/internal/cache"
	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/stretchr/testify/require"
)

func TestMemoryCache(t *testing.T) {
	c := cache.NewMemory()
	ctx := context.Background()
	_, hit, err := c.Get(ctx)
	require.NoError(t, err)
	require.False(t, hit)

	require.NoError(t, c.Set(ctx, model.Portfolio{Profile: model.Profile{FullName: "Mohd Hujaifa"}}))
	p, hit, err := c.Get(ctx)
	require.NoError(t, err)
	require.True(t, hit)
	require.Equal(t, "Mohd Hujaifa", p.Profile.FullName)

	require.NoError(t, c.Invalidate(ctx))
	_, hit, err = c.Get(ctx)
	require.NoError(t, err)
	require.False(t, hit)
}
