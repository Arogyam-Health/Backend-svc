package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerCallsFunctions(t *testing.T) {
	var refreshCalled int32

	refreshFn := func() {
		atomic.AddInt32(&refreshCalled, 1)
	}

	ctx := context.Background()
	StartTokenRefresh(ctx, refreshFn)

	time.Sleep(2 * time.Second)

	if atomic.LoadInt32(&refreshCalled) == 0 {
		t.Fatal("refreshFn was not called")
	}
}
