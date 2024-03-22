package Flood

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	// Create a new Flood instance with a count of 5 and a timeDelta of 10 seconds
	f := NewFlood(5, 10)
	wg := sync.WaitGroup{}

	// Test case 1: User has made less than the max allowed calls in the last timeDelta seconds
	t.Run("User has made less than the max allowed calls", func(t *testing.T) {
		ctx := context.Background()
		userID := int64(1)

		// Add a few requests for the user
		f.Add(userID, time.Now().Add(-5*time.Second))
		f.Add(userID, time.Now().Add(-7*time.Second))

		// Check if the user is allowed to perform an action
		ch, err := f.Check(ctx, userID)

		assert.NoError(t, err)
		assert.True(t, ch)
	})

	// Test case 2: User has made exactly the max allowed calls in the last timeDelta seconds
	t.Run("User has made exactly the max allowed calls", func(t *testing.T) {
		ctx := context.Background()
		userID := int64(2)

		// Add the max allowed calls for the user
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				f.Add(userID, time.Now().Add(-time.Duration(i)*time.Second))
			}(i)
		}
		wg.Wait()

		// Check if the user is allowed to perform an action
		ch, err := f.Check(ctx, userID)

		assert.NoError(t, err)
		assert.True(t, ch)
	})

	// Test case 3: User has made more than the max allowed calls in the last timeDelta seconds
	t.Run("User has made more than the max allowed calls", func(t *testing.T) {
		ctx := context.Background()
		userID := int64(3)

		// Add more than the max allowed calls for the user
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				f.Add(userID, time.Now().Add(-time.Duration(i)*time.Second))
			}(i)
		}
		wg.Wait()

		// Check if the user is allowed to perform an action
		ch, err := f.Check(ctx, userID)

		assert.NoError(t, err)
		assert.False(t, ch)
	})

	// Test case 4: Context is canceled before the check
	t.Run("Context is canceled before the check", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		userID := int64(4)

		// Cancel the context before the check
		cancel()

		// Check if the user is allowed to perform an action
		_, err := f.Check(ctx, userID)
		assert.ErrorContains(t, err, "context canceled")
		// context deadline exceeded
	})

	// Test case 4: Context deadline exceeded
	t.Run("Context deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		userID := int64(1)

		time.Sleep(2 * time.Second)
		_, err := f.Check(ctx, userID)

		assert.ErrorContains(t, err, "context deadline exceeded")
	})
}
