package p

import (
	"context"
	"testing"
)

// This test doesn't run yet.
func TestUpdateDB(t *testing.T) {
	m := PubSubMessage{}
	ctx := context.Background()

	err := UpdateDB(ctx, m)
	if err != nil {
		t.Fatalf("faild test : %v", err)
	}
}
