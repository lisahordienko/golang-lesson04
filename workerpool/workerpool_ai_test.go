package workerpool

import (
	"context"
	"testing"
	"time"
)

func TestRunPool_TimeoutMechanism(t *testing.T) {
	jobs := make(chan Job, 1)
	jobs <- Job{
		ID: "slow-ai-job",
		Fetch: func(ctx context.Context) (int, error) {
			select {
			case <-time.After(2 * time.Second):
				return 999, nil
			case <-ctx.Done():
				return 0, ctx.Err()
			}
		},
	}
	close(jobs)

	start := time.Now()
	results := RunPool(jobs, 1, 150*time.Millisecond)

	var got Result
	for r := range results {
		got = r
	}

	if got.Err == nil {
		t.Fatal("очікувалась помилка тайм-ауту, отримано nil")
	}
	if got.Err != context.DeadlineExceeded {
		t.Fatalf("Err = %v, want %v", got.Err, context.DeadlineExceeded)
	}
	if got.JobID != "slow-ai-job" {
		t.Fatalf("JobID = %q, want %q", got.JobID, "slow-ai-job")
	}

	elapsed := time.Since(start)
	if elapsed > 700*time.Millisecond {
		t.Fatalf("тайм-аут спрацював пізно: %v", elapsed)
	}
}
