package Flood

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"
)

type Flood struct {
	users   map[int64][]time.Time
	count   int64
	timeDur time.Duration
	mu      sync.RWMutex
}

// Add new time for given userID and return slice of times
func (f *Flood) Add(userID int64, t time.Time) []time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.users[userID]; !ok {
		f.users[userID] = []time.Time{t}
	} else {
		f.users[userID] = append(f.users[userID], t)
	}
	return f.users[userID]
}

// How many times has the function been called in the last f.timeDur seconds
func (f *Flood) Count(ts []time.Time, t time.Time) int {
	i, _ := slices.BinarySearchFunc(ts, t.Add(-f.timeDur*time.Second), func(a, b time.Time) int {
		if a == b {
			return 0
		} else if a.Before(b) {
			return -1
		} else {
			return 1
		}
	})
	return len(ts) - i
}

// Returns a new Flood
func NewFlood(c int64, td time.Duration) Flood {
	return Flood{
		users:   make(map[int64][]time.Time),
		count:   c,
		timeDur: td,
	}
}

// Check checks if a user is allowed to perform an action based on the count
// of requests in the last timeDur seconds.
func (f *Flood) Check(ctx context.Context, userID int64) (ch bool, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()
	now := time.Now()

	// Add the current time for the user to the slice of times for the user
	ts := f.Add(userID, now)
	// Sort the slice of times for the user by time
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Before(ts[j])
	})
	// Count how many times the user has made a request in the last timeDur seconds
	count := f.Count(ts, now)

	// count is greater than the max allowed calls
	if int64(count) > f.count {
		ch = false
	} else {
		ch = true
	}
	return
}
