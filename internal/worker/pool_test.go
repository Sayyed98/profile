package worker_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mohdhujaifa/profile/internal/worker"
	"github.com/stretchr/testify/require"
)

func TestPoolRunsJobs(t *testing.T) {
	ctx := context.Background()
	p := worker.NewPool(ctx, 2)
	defer p.Stop()

	var n atomic.Int32
	require.True(t, p.Submit(func(context.Context) { n.Add(1) }))
	require.Eventually(t, func() bool { return n.Load() == 1 }, time.Second, 10*time.Millisecond)
}
