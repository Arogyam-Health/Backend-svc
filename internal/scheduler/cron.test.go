package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerCallsFunctions(t *testing.T) {
	var syncCalled int32
	var refreshCalled int32

	syncFn := func() {
		atomic.AddInt32(&syncCalled, 1)
	}

	refreshFn := func() {
		atomic.AddInt32(&refreshCalled, 1)
	}

	ctx := context.Background()
	Start(ctx, syncFn, refreshFn)

	time.Sleep(2 * time.Second)

	if atomic.LoadInt32(&syncCalled) == 0 {
		t.Fatal("syncFn was not called")
	}
}
