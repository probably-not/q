package pico

import (
	"math/rand"
	"testing"
)

func TestPush(t *testing.T) {
	randomAmountOfJobs := rand.Intn(62)

	testCases := []struct {
		desc               string
		queue              Q
		expectedIsFull     bool
		nextExpectedIsFull bool
		expectedIdx        int
		nextExpectedIdx    int
		queueSizeFactor    int
	}{
		{
			desc:               "Zero value of queue allows pushing",
			queue:              NewQ(),
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        0,
			nextExpectedIdx:    1,
			queueSizeFactor:    6,
		},
		{
			desc:               "Queue with random amount of jobs less than size factor allows pushing",
			queue:              Q(randomAmountOfJobs),
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        randomAmountOfJobs,
			nextExpectedIdx:    randomAmountOfJobs + 1,
			queueSizeFactor:    6,
		},
		{
			desc:               "Queue at the size factor is marked as full and cannot be pushed to",
			queue:              Q(63),
			expectedIsFull:     true,
			nextExpectedIsFull: false,
			expectedIdx:        -1,
			nextExpectedIdx:    0,
			queueSizeFactor:    6,
		},
		{
			desc:               "Overflow is protected against",
			queue:              Q(4294967295),
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        63,
			nextExpectedIdx:    0,
			queueSizeFactor:    6,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(subT *testing.T) {
			idx, isFull := tC.queue.Push(tC.queueSizeFactor)
			if tC.expectedIdx != idx {
				subT.Errorf("expected the returned index to be %d, got %d", tC.expectedIdx, idx)
			}

			if tC.expectedIsFull != isFull {
				subT.Errorf("expected isFull to be %t, got %t", tC.expectedIsFull, isFull)
			}

			tC.queue.PushCommit()

			idx, isFull = tC.queue.Push(tC.queueSizeFactor)
			if tC.nextExpectedIdx != idx {
				subT.Errorf("expected the next returned index to be %d, got %d", tC.nextExpectedIdx, idx)
			}

			if tC.nextExpectedIsFull != isFull {
				subT.Errorf("expected next isFull to be %t, got %t", tC.nextExpectedIsFull, isFull)
			}
		})
	}
}
